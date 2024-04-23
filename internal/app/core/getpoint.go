package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

func (s *pickUpPointCoreService) GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error) {
	point, err := s.cache.GetPoint(id)
	if err == nil {
		return point, nil
	}
	return s.pointService.GetPoint(ctx, id)
}
