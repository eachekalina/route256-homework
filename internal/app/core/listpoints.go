package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

func (s *pickUpPointCoreService) ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error) {
	return s.pointService.ListPoints(ctx)
}
