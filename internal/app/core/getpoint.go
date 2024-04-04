package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

func (s *pickUpPointCoreService) GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error) {
	return s.pointService.GetPoint(ctx, id)
}
