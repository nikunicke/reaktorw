package pipeline

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
)

type fifo struct {
	proc Processor
}

func FIFO(proc Processor) StageRunner {
	return fifo{proc: proc}
}

func (r fifo) Run(ctx context.Context, params StageParams) {
	for {
		select {
		case <-ctx.Done():
			return
		case payloadIn, ok := <-params.Input():
			if !ok {
				return
			}
			payloadOut, err := r.proc.Process(ctx, payloadIn)
			if err != nil {
				newErr := xerrors.New("pipeline error")
				maybeEmitError(newErr, params.Error())
				return
			}
			if payloadOut == nil {
				payloadIn.MarkAsProcessed()
				continue
			}
			select {
			case params.Output() <- payloadOut:
			case <-ctx.Done():
				return
			}
		}
	}
}

type fixedWorkerPool struct {
	fifos []StageRunner
}

func FixedWorkerPool(proc Processor, numWorkers uint) StageRunner {
	if numWorkers <= 0 {
		panic("pipeline: FixedWorkerPool numWorkers must be > 0")
	}
	fifos := make([]StageRunner, int(numWorkers))
	for i := 0; i < int(numWorkers); i++ {
		fifos[i] = FIFO(proc)
	}
	return &fixedWorkerPool{fifos: fifos}
}

func (r *fixedWorkerPool) Run(ctx context.Context, params StageParams) {
	var wg sync.WaitGroup

	for i := 0; i < len(r.fifos); i++ {
		wg.Add(1)
		go func(fifoIndex int) {
			r.fifos[fifoIndex].Run(ctx, params)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
