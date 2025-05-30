package cqrx

import "github.com/ThreeDotsLabs/watermill/components/cqrs"

func DefaultMarshaler() cqrs.CommandEventMarshaler {
	return &cqrs.JSONMarshaler{
		GenerateName: cqrs.NamedStruct(cqrs.StructName),
	}
}
