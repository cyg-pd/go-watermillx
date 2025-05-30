package metrics

var (
	handlerLabelKeys = []string{
		labelKeyHandlerName,
	}

	// defaultHandlerExecutionTimeBuckets are one order of magnitude smaller than default buckets (5ms~10s),
	// because the handler execution times are typically shorter (µs~ms range).
	defaultHandlerExecutionTimeBuckets = []float64{
		0.0005,
		0.001,
		0.0025,
		0.005,
		0.01,
		0.025,
		0.05,
		0.1,
		0.25,
		0.5,
		1,
	}
)
