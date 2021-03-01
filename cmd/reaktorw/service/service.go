package service

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
)

// Service represent a service for the reaktor warehouse application
type Service interface {
	Name() string
	Run(context.Context) error
}

// Group is a list of service instances
type Group []Service

// Run executes all service instances in the group
func (g Group) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	runCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	var wg sync.WaitGroup
	errCh := make(chan error, len(g))
	wg.Add(len(g))

	for _, service := range g {
		go func(s Service) {
			defer wg.Done()
			if err := s.Run(runCtx); err != nil {
				errCh <- xerrors.New("service fail")
				cancelFn()
			}
		}(service)
	}
	<-runCtx.Done()
	wg.Wait()

	var err error
	close(errCh)
	for serviceErr := range errCh {
		// err = multierror.Append(err, serviceErr)
		err = serviceErr
	}
	return err
}
