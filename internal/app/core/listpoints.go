package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

func (s *pickUpPointCoreService) ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error) {
	ctx, span := s.tracer.Start(ctx, "ListPoints")
	defer span.End()

	points, err := s.redis.GetPointList(ctx)
	if err == nil {
		return points, nil
	}
	s.log.Log(err.Error())
	points, err = s.pointService.ListPoints(ctx)
	if err != nil {
		return nil, err
	}
	err = s.redis.SetPointList(ctx, points)
	if err != nil {
		return nil, err
	}
	return points, nil
}
