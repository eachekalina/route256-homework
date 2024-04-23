package db

import "context"

type Dummy struct {
}

func (d Dummy) RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error {
	return f(ctx)
}
