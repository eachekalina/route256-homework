package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

type CreatePointRequest struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func (s *pickUpPointCoreService) CreatePoint(ctx context.Context, req CreatePointRequest) error {
	ctx, span := s.tracer.Start(ctx, "CreatePoint")
	defer span.End()

	point := pickuppoint.PickUpPoint{
		Id:      req.Id,
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}
	err := s.pointService.CreatePoint(ctx, point)
	if err != nil {
		return err
	}
	s.cache.PutPoint(ctx, point)
	return nil
}
