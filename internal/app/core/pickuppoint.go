//go:generate mockgen -source=./pickuppoint.go -destination=./mocks/pickuppoint.go -package=mocks

package core

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"homework/internal/app/logger"
	"homework/internal/app/pickuppoint"
)

type PickUpPointCoreService interface {
	CreatePoint(ctx context.Context, req CreatePointRequest) error
	ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error)
	GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error)
	UpdatePoint(ctx context.Context, req UpdatePointRequest) error
	DeletePoint(ctx context.Context, id uint64) error
	SetCache(cache Cache)
	SetRedis(redis Redis)
}

type Cache interface {
	PutPoint(ctx context.Context, point pickuppoint.PickUpPoint)
	GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error)
	DeletePoint(ctx context.Context, id uint64)
}

type Redis interface {
	GetPointList(ctx context.Context) ([]pickuppoint.PickUpPoint, error)
	SetPointList(ctx context.Context, points []pickuppoint.PickUpPoint) error
}

type pickUpPointCoreService struct {
	pointService PickUpPointService
	log          logger.Logger
	cache        Cache
	redis        Redis
	tracer       trace.Tracer
}

type PickUpPointService interface {
	CreatePoint(ctx context.Context, point pickuppoint.PickUpPoint) error
	ListPoints(ctx context.Context) ([]pickuppoint.PickUpPoint, error)
	GetPoint(ctx context.Context, id uint64) (pickuppoint.PickUpPoint, error)
	UpdatePoint(ctx context.Context, point pickuppoint.PickUpPoint) error
	DeletePoint(ctx context.Context, id uint64) error
}

func NewPickUpPointCoreService(pointService PickUpPointService, log logger.Logger) PickUpPointCoreService {
	return &pickUpPointCoreService{
		pointService: pointService,
		log:          log,
		cache:        NilCache{},
		redis:        NilRedis{},
		tracer:       otel.Tracer("internal/app/core/pickuppoint"),
	}
}

func (s *pickUpPointCoreService) SetCache(cache Cache) {
	s.cache = cache
}

func (s *pickUpPointCoreService) SetRedis(redis Redis) {
	s.redis = redis
}
