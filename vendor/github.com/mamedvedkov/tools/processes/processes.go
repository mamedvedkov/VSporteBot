package processes

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Process func(ctx context.Context) (err error)

func RunAndWait(ctx context.Context, fn func() error, done func()) error {
	errCh := make(chan error)

	go func() {
		errCh <- fn()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		done()

		return ctx.Err()
	}
}

func RunParallelAndWait(ctx context.Context, processes ...Process) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, fn := range processes {
		internalFn := fn

		group.Go(func() error {
			return internalFn(ctx)
		})
	}

	return group.Wait()
}
