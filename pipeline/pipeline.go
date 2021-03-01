package pipeline

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
)

var _ StageParams = (*workerParams)(nil)

type workerParams struct {
	stage int

	inCh  <-chan Payload
	outCh chan<- Payload
	errCh chan<- error
}

func (p *workerParams) StageIndex() int        { return p.stage }
func (p *workerParams) Input() <-chan Payload  { return p.inCh }
func (p *workerParams) Output() chan<- Payload { return p.outCh }
func (p *workerParams) Error() chan<- error    { return p.errCh }

type Pipeline struct {
	stages []StageRunner
}

func New(stages ...StageRunner) *Pipeline {
	return &Pipeline{stages: stages}
}

// Process
func (p *Pipeline) Process(ctx context.Context, source Source, sink Sink) error {
	var wg sync.WaitGroup
	pCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	stageCh := make([]chan Payload, len(p.stages)+1)
	errCh := make(chan error, len(p.stages)+2)
	for i := 0; i < len(p.stages)+1; i++ {
		stageCh[i] = make(chan Payload)
	}
	for i := 0; i < len(p.stages); i++ {
		wg.Add(1)
		go func(stageIndex int) {
			p.stages[stageIndex].Run(pCtx, &workerParams{
				stage: stageIndex,
				inCh:  stageCh[stageIndex],
				outCh: stageCh[stageIndex+1],
				errCh: errCh,
			})
			close(stageCh[stageIndex+1])
			wg.Done()
		}(i)
	}
	wg.Add(2)
	go func() {
		sourceWorker(pCtx, source, stageCh[0], errCh)
		close(stageCh[0])
		wg.Done()
	}()
	go func() {
		sinkWorker(pCtx, sink, stageCh[len(stageCh)-1], errCh)
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(errCh)
		cancelFn()
	}()

	var errAll error

	for err := range errCh {
		// errAll = multierror.Append(errAll, err)
		errAll = err
		cancelFn()
	}
	return errAll
}

// sourceWorker
func sourceWorker(ctx context.Context, source Source, outCh chan<- Payload, errCh chan<- error) {
	for source.Next(ctx) {
		payload := source.Payload()
		select {
		case outCh <- payload:
		case <-ctx.Done():
			return
		}
	}
	if err := source.Error(); err != nil {
		newErr := xerrors.New("source error")
		maybeEmitError(newErr, errCh)
	}
}

func sinkWorker(ctx context.Context, sink Sink, inCh <-chan Payload, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case payload, ok := <-inCh:
			if !ok {
				return
			}
			if err := sink.Consume(ctx, payload); err != nil {
				newErr := xerrors.New("sink error")
				maybeEmitError(newErr, errCh)
				return
			}
			payload.MarkAsProcessed()
		}
	}
}

func maybeEmitError(err error, errCh chan<- error) {
	select {
	case errCh <- err:
	default:
	}
}
