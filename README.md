# go-watermillx

[![tag](https://img.shields.io/github/tag/cyg-pd/go-watermillx.svg)](https://github.com/cyg-pd/go-watermillx/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/cyg-pd/go-watermillx?status.svg)](https://pkg.go.dev/github.com/cyg-pd/go-watermillx)
![Build Status](https://github.com/cyg-pd/go-watermillx/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/cyg-pd/go-watermillx)](https://goreportcard.com/report/github.com/cyg-pd/go-watermillx)
[![Coverage](https://img.shields.io/codecov/c/github/cyg-pd/go-watermillx)](https://codecov.io/gh/cyg-pd/go-watermillx)
[![Contributors](https://img.shields.io/github/contributors/cyg-pd/go-watermillx)](https://github.com/cyg-pd/go-watermillx/graphs/contributors)
[![License](https://img.shields.io/github/license/cyg-pd/go-watermillx)](./LICENSE)

## 🚀 Install

```sh
go get github.com/cyg-pd/go-watermillx@v1
```

This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0.

## 💡 Usage

You can import `watermillx` using:

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/cyg-pd/go-otelx"
	_ "github.com/cyg-pd/go-otelx/autoconf"
	"github.com/cyg-pd/go-watermillx/cqrx"
	"github.com/cyg-pd/go-watermillx/driver"
	_ "github.com/cyg-pd/go-watermillx/driver/kafka"
	"github.com/cyg-pd/go-watermillx/manager"
	"github.com/cyg-pd/go-watermillx/opentelemetry/metrics"
	"github.com/cyg-pd/go-watermillx/opentelemetry/trace"
	"github.com/cyg-pd/go-watermillx/router"
)

var tracer = otelx.Tracer()

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	r := must(router.New())
	r.AddMiddleware(middleware.Recoverer)
	r.AddMiddleware(trace.HandlerMiddleware())
	r.AddPublisherDecorators(trace.PublisherDecorator())

	mb := metrics.NewOpenTelemetryMetricsBuilder(otelx.Meter(), "", "")
	mb.AttachRouter(r)

	da := must(driver.New("kafka", `{"Brokers": ["localhost:9092"], "InitializeTopicDetails":{"NumPartitions": 3}}`))
	db := must(driver.New("kafka", `{"Brokers": ["localhost:9092"], "InitializeTopicDetails":{"NumPartitions": 10}}`))

	m := manager.New(r)
	m.Add("kafka-a", da)
	m.Add("kafka-b", db)

	cqrsA := m.MustUseCQRS("kafka-a")
	cqrsA.AddPublisherDecorator(trace.PublisherDecorator())
	cqrsA.AddPublisherDecorator(mb.DecoratePublisher)
	cqrsA.AddSubscriberDecorator(mb.DecorateSubscriber)

	cqrsB := m.MustUseCQRS("kafka-b")
	cqrsB.AddPublisherDecorator(trace.PublisherDecorator())
	cqrsB.AddPublisherDecorator(mb.DecoratePublisher)
	cqrsB.AddSubscriberDecorator(mb.DecorateSubscriber)

	commandBus := cqrsA.MustCommandBus(cqrx.WithCommandBusOnSend(commandBusOnSendHook))
	commandProcessor := cqrsA.MustCommandProcessor(cqrx.WithCommandProcessorOnHandle(commandProcessorOnHandleHook))

	eventBus := cqrsB.MustEventBus(cqrx.WithEventBusOnPublish(eventBusOnPublish))
	eventProcessor := cqrsB.MustEventProcessor(cqrx.WithEventProcessorOnHandle(eventProcessorOnHandle))

	must0(
		commandProcessor.AddHandlers(
			cqrs.NewCommandHandler("BookRoomHandler", BookRoomHandler{eventBus}.Handle),
			cqrs.NewCommandHandler("OrderBeerHandler", OrderBeerHandler{eventBus}.Handle),
		),
	)
	must0(
		eventProcessor.AddHandlers(
			cqrs.NewEventHandler(
				"OrderBeerOnRoomBooked",
				OrderBeerOnRoomBooked{commandBus}.Handle,
			),
			cqrs.NewEventHandler(
				"LogBeerOrdered",
				func(ctx context.Context, event *BeerOrdered) error {
					slog.Info("Beer ordered", slog.String("room_id", event.RoomId))
					return nil
				},
			),
			cqrs.NewEventHandler(
				"BookingsFinancialReport",
				NewBookingsFinancialReport().Handle,
			),
		),
	)

	// publish BookRoom commands every second to simulate incoming traffic
	go publishCommands(commandBus)

	// processors are based on router, so they will work when router will start
	if err := r.Run(context.Background()); err != nil {
		panic(err)
	}
}

func publishCommands(commandBus *cqrs.CommandBus) func() {
	i := 0
	for {
		i++
		func() {
			ctx, span := tracer.Start(context.Background(), "publishCommands")
			defer span.End()

			startDate := time.Now()
			endDate := time.Now().Add(time.Hour * 24 * 3)

			bookRoomCmd := &BookRoom{
				RoomId:    fmt.Sprintf("%d", i),
				GuestName: "John",
				StartDate: startDate,
				EndDate:   endDate,
			}
			if err := commandBus.Send(ctx, bookRoomCmd); err != nil {
				panic(err)
			}

			time.Sleep(time.Second)
		}()
	}
}

func commandBusOnSendHook(params cqrs.CommandBusOnSendParams) error {
	slog.Info(
		"Sending command",
		slog.String("command_name", params.CommandName),
	)
	return nil
}

func commandProcessorOnHandleHook(params cqrs.CommandProcessorOnHandleParams) error {
	start := time.Now()
	err := params.Handler.Handle(params.Message.Context(), params.Command)
	slog.Info(
		"Command handled",
		slog.String("command_name", params.CommandName),
		slog.Duration("duration", time.Since(start)),
		slog.Any("err", err),
	)
	return err
}

func eventBusOnPublish(params cqrs.OnEventSendParams) error {
	slog.Info("Publishing event", slog.String("event_name", params.EventName))
	return nil
}

func eventProcessorOnHandle(params cqrs.EventProcessorOnHandleParams) error {
	start := time.Now()
	err := params.Handler.Handle(params.Message.Context(), params.Event)
	slog.Info(
		"Event handled",
		slog.String("event_name", params.EventName),
		slog.Duration("duration", time.Since(start)),
		slog.Any("err", err),
	)
	return err
}

type BookRoom struct {
	RoomId    string    `json:"room_id,omitempty"`
	GuestName string    `json:"guest_name,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

type RoomBooked struct {
	ReservationId string    `json:"reservation_id,omitempty"`
	RoomId        string    `json:"room_id,omitempty"`
	GuestName     string    `json:"guest_name,omitempty"`
	Price         int64     `json:"price,omitempty"`
	StartDate     time.Time `json:"start_date,omitempty"`
	EndDate       time.Time `json:"end_date,omitempty"`
}

type OrderBeer struct {
	RoomId string `json:"room_id,omitempty"`
	Count  int64  `json:"count,omitempty"`
}

type BeerOrdered struct {
	RoomId string `json:"room_id,omitempty"`
	Count  int64  `json:"count,omitempty"`
}

// BookRoomHandler is a command handler, which handles BookRoom command and emits RoomBooked.
//
// In CQRS, one command must be handled by only one handler.
// When another handler with this command is added to command processor, error will be returned.
type BookRoomHandler struct {
	eventBus *cqrs.EventBus
}

func (b BookRoomHandler) Handle(ctx context.Context, cmd *BookRoom) error {
	ctx, span := tracer.Start(ctx, "BookRoomHandler")
	defer span.End()

	// some random price, in production you probably will calculate in wiser way
	price := (rand.Int63n(40) + 1) * 10

	slog.Info(
		"Booked room",
		"room_id", cmd.RoomId,
		"guest_name", cmd.GuestName,
		"start_date", cmd.StartDate,
		"end_date", cmd.EndDate,
	)

	// RoomBooked will be handled by OrderBeerOnRoomBooked event handler,
	// in future RoomBooked may be handled by multiple event handler
	if err := b.eventBus.Publish(ctx, &RoomBooked{
		ReservationId: watermill.NewUUID(),
		RoomId:        cmd.RoomId,
		GuestName:     cmd.GuestName,
		Price:         price,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
	}); err != nil {
		return err
	}

	return nil
}

// OrderBeerOnRoomBooked is an event handler, which handles RoomBooked event and emits OrderBeer command.
type OrderBeerOnRoomBooked struct {
	commandBus *cqrs.CommandBus
}

func (o OrderBeerOnRoomBooked) Handle(ctx context.Context, event *RoomBooked) error {
	ctx, span := tracer.Start(ctx, "OrderBeerOnRoomBooked")
	defer span.End()

	orderBeerCmd := &OrderBeer{
		RoomId: event.RoomId,
		Count:  rand.Int63n(10) + 1,
	}

	return o.commandBus.Send(ctx, orderBeerCmd)
}

// OrderBeerHandler is a command handler, which handles OrderBeer command and emits BeerOrdered.
// BeerOrdered is not handled by any event handler, but we may use persistent Pub/Sub to handle it in the future.
type OrderBeerHandler struct {
	eventBus *cqrs.EventBus
}

func (o OrderBeerHandler) Handle(ctx context.Context, cmd *OrderBeer) error {
	ctx, span := tracer.Start(ctx, "OrderBeerHandler")
	defer span.End()

	if rand.Int63n(10) == 0 {
		// sometimes there is no beer left, command will be retried
		return fmt.Errorf("no beer left for room %s, please try later", cmd.RoomId)
	}

	if err := o.eventBus.Publish(ctx, &BeerOrdered{
		RoomId: cmd.RoomId,
		Count:  cmd.Count,
	}); err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("%d beers ordered to room %s", cmd.Count, cmd.RoomId))
	return nil
}

// BookingsFinancialReport is a read model, which calculates how much money we may earn from bookings.
// Like OrderBeerOnRoomBooked, it listens for RoomBooked event.
//
// This implementation is just writing to the memory. In production, you will probably will use some persistent storage.
type BookingsFinancialReport struct {
	handledBookings map[string]struct{}
	totalCharge     int64
	lock            sync.Mutex
}

func NewBookingsFinancialReport() *BookingsFinancialReport {
	return &BookingsFinancialReport{handledBookings: map[string]struct{}{}}
}

func (b *BookingsFinancialReport) Handle(ctx context.Context, event *RoomBooked) error {
	_, span := tracer.Start(ctx, "BookingsFinancialReport")
	defer span.End()

	// Handle may be called concurrently, so it need to be thread safe.
	b.lock.Lock()
	defer b.lock.Unlock()

	// When we are using Pub/Sub which doesn't provide exactly-once delivery semantics, we need to deduplicate messages.
	// GoChannel Pub/Sub provides exactly-once delivery,
	// but let's make this example ready for other Pub/Sub implementations.
	if _, ok := b.handledBookings[event.ReservationId]; ok {
		return nil
	}
	b.handledBookings[event.ReservationId] = struct{}{}

	b.totalCharge += event.Price

	slog.Info(fmt.Sprintf(">>> Already booked rooms for $%d\n", b.totalCharge))
	return nil
}
```
