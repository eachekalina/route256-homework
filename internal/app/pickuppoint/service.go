package pickuppoint

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	Create(ctx context.Context, point PickUpPoint) error
	List(ctx context.Context) ([]PickUpPoint, error)
	Get(ctx context.Context, id uint64) (PickUpPoint, error)
	Update(ctx context.Context, point PickUpPoint) error
	Delete(ctx context.Context, id uint64) error
}

type TransactionManager interface {
	RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error
}

var ErrIdAlreadyExists = errors.New("item with such id already exists")
var ErrNoItemFound = errors.New("no such item found")

// Service allows concurrent working on pick-up points.
type Service struct {
	repo   Repository
	tm     TransactionManager
	tracer trace.Tracer
}

// NewService creates a new Service.
func NewService(repo Repository, tm TransactionManager) *Service {
	return &Service{
		repo:   repo,
		tm:     tm,
		tracer: otel.Tracer("internal/app/pickuppoint/service"),
	}
}

// CreatePoint creates a pick-up point.
func (s *Service) CreatePoint(ctx context.Context, point PickUpPoint) error {
	ctx, span := s.tracer.Start(ctx, "CreatePoint")
	defer span.End()

	return s.tm.RunSerializable(ctx, func(ctxTX context.Context) error {
		return s.repo.Create(ctxTX, point)
	})
}

// ListPoints prints a slice of all pick-up points.
func (s *Service) ListPoints(ctx context.Context) ([]PickUpPoint, error) {
	ctx, span := s.tracer.Start(ctx, "ListPoints")
	defer span.End()

	var points []PickUpPoint
	err := s.tm.RunSerializable(ctx, func(ctxTX context.Context) error {
		var err error
		points, err = s.repo.List(ctxTX)
		return err
	})
	return points, err
}

// GetPoint prints a specified pick-up point.
func (s *Service) GetPoint(ctx context.Context, id uint64) (PickUpPoint, error) {
	ctx, span := s.tracer.Start(ctx, "GetPoint")
	defer span.End()

	var point PickUpPoint
	err := s.tm.RunSerializable(ctx, func(ctxTX context.Context) error {
		var err error
		point, err = s.repo.Get(ctxTX, id)
		return err
	})
	return point, err
}

// UpdatePoint updates a pick-up point info.
func (s *Service) UpdatePoint(ctx context.Context, point PickUpPoint) error {
	ctx, span := s.tracer.Start(ctx, "UpdatePoint")
	defer span.End()

	return s.tm.RunSerializable(ctx, func(ctxTX context.Context) error {
		return s.repo.Update(ctxTX, point)
	})
}

// DeletePoint deletes a pick-up point.
func (s *Service) DeletePoint(ctx context.Context, id uint64) error {
	ctx, span := s.tracer.Start(ctx, "DeletePoint")
	defer span.End()

	return s.tm.RunSerializable(ctx, func(ctxTX context.Context) error {
		return s.repo.Delete(ctxTX, id)
	})
}
