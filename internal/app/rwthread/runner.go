package rwthread

import (
	"context"
	"golang.org/x/sync/errgroup"
	"homework/internal/app/logger"
)

type RunFunc func(ctx context.Context) error

type Runner struct {
	log   logger.Logger
	read  chan RunFunc
	write chan RunFunc
}

func NewRunner(log logger.Logger) *Runner {
	return &Runner{log: log}
}

func (r *Runner) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return r.writeThread(ctx)
	})

	eg.Go(func() error {
		return r.readThread(ctx)
	})

	return eg.Wait()
}

func (r *Runner) writeThread(ctx context.Context) error {
	r.write = make(chan RunFunc)
	for {
		r.log.Log("write thread: waiting for request")
		select {
		case <-ctx.Done():
			r.log.Log("write thread: closing")
			close(r.write)
			return ctx.Err()
		case f := <-r.write:
			err := f(ctx)
			if err != nil {
				r.log.Log("write thread: error: %v", err)
			}
		}
	}
}

func (r *Runner) readThread(ctx context.Context) error {
	r.read = make(chan RunFunc)

	for {
		r.log.Log("read thread: waiting for request")
		select {
		case <-ctx.Done():
			r.log.Log("read thread: closing")
			close(r.read)
			return ctx.Err()
		case f := <-r.read:
			err := f(ctx)
			if err != nil {
				r.log.Log("read thread: error: %v", err)
			}
		}
	}
}

func (r *Runner) RunRead(f RunFunc) {
	r.read <- f
}

func (r *Runner) RunWrite(f RunFunc) {
	r.write <- f
}
