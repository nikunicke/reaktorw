package pipeline

import "context"

type Payload interface {
	Clone() Payload
	MarkAsProcessed()
}

type Processor interface {
	Process(context.Context, Payload) (Payload, error)
}

type ProcessorFunc func(context.Context, Payload) (Payload, error)

func (f ProcessorFunc) Process(ctx context.Context, p Payload) (Payload, error) {
	return f(ctx, p)
}

type StageRunner interface {
	Run(context.Context, StageParams)
}

type StageParams interface {
	StageIndex() int
	Input() <-chan Payload
	Output() chan<- Payload
	Error() chan<- error
}

type Source interface {
	Next(context.Context) bool
	Payload() Payload
	Error() error
}

type Sink interface {
	Consume(context.Context, Payload) error
}
