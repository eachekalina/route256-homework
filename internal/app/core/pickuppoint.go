//go:generate mockgen -source=./pickuppoint.go -destination=./mocks/pickuppoint.go -package=mocks

package core

import (
	"context"
	"homework/internal/app/pickuppoint"
)

type PickUpPointCoreService interface {
	CreatePoint(ctx context.Context, req CreatePointRequest) error
	ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error)
	GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error)
	UpdatePoint(ctx context.Context, req UpdatePointRequest) error
	DeletePoint(ctx context.Context, id uint64) error
}

type pickUpPointCoreService struct {
	pointService PickUpPointService
}

type PickUpPointService interface {
	CreatePoint(ctx context.Context, point pickuppoint.PickUpPoint) error
	ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error)
	GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error)
	UpdatePoint(ctx context.Context, point pickuppoint.PickUpPoint) error
	DeletePoint(ctx context.Context, id uint64) error
}

func NewPickUpPointCoreService(pointService PickUpPointService) PickUpPointCoreService {
	return &pickUpPointCoreService{pointService: pointService}
}
