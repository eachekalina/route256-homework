package core

import (
	"context"
	"errors"
	"homework/internal/app/pickuppoint"
)

type NilCache struct {
}

func (n NilCache) PutPoint(ctx context.Context, point pickuppoint.PickUpPoint) {
}

func (n NilCache) GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error) {
	return pickuppoint.PickUpPoint{}, errors.New("not found")
}

func (n NilCache) DeletePoint(ctx context.Context, id uint64) {
}
