package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

type UpdatePointRequest struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func (s *pickUpPointCoreService) UpdatePoint(ctx context.Context, req UpdatePointRequest) error {
	point := pickuppoint.PickUpPoint{
		Id:      req.Id,
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
	}
	err := s.pointService.UpdatePoint(ctx, point)
	if err != nil {
		return err
	}
	s.cache.PutPoint(point)
	return nil
}
