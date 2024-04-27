package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

func (s *pickUpPointCoreService) GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error) {
	ctx, span := s.tracer.Start(ctx, "GetPoint")
	defer span.End()

	point, err := s.cache.GetPoint(ctx, id)
	if err == nil {
		return point, nil
	}
	return s.pointService.GetPoint(ctx, id)
}
