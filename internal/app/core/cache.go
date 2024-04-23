package core

import (
	"errors"
	"homework/internal/app/pickuppoint"
)

type NilCache struct {
}

func (n NilCache) PutPoint(point pickuppoint.PickUpPoint) {
}

func (n NilCache) GetPoint(id uint64) (pickuppoint.PickUpPoint, error) {
	return pickuppoint.PickUpPoint{}, errors.New("not found")
}

func (n NilCache) DeletePoint(id uint64) {
}
