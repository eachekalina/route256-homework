package core

import (
	"context"
	"errors"
	"homework/internal/app/pickuppoint"
)

type NilRedis struct {
}

func (n NilRedis) GetPointList(ctx context.Context) ([]pickuppoint.PickUpPoint, error) {
	return nil, errors.New("not found")
}

func (n NilRedis) SetPointList(ctx context.Context, points []pickuppoint.PickUpPoint) error {
	return nil
}
