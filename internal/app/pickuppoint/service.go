package pickuppoint

import (
	"context"
	"errors"
)

type Repository interface {
	Create(ctx context.Context, point PickUpPoint) error
	List(ctx context.Context) ([]PickUpPoint, error)
	Get(ctx context.Context, id uint64) (PickUpPoint, error)
	Update(ctx context.Context, point PickUpPoint) error
	Delete(ctx context.Context, id uint64) error
}

var ErrIdAlreadyExists = errors.New("item with such id already exists")
var ErrNoItemFound = errors.New("no such item found")

// Service allows concurrent working on pick-up points.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreatePoint creates a pick-up point.
func (s *Service) CreatePoint(ctx context.Context, point PickUpPoint) error {
	return s.repo.Create(ctx, point)
}

// ListPoints prints a slice of all pick-up points.
func (s *Service) ListPoints(ctx context.Context) ([]PickUpPoint, error) {
	return s.repo.List(ctx)
}

// GetPoint prints a specified pick-up point.
func (s *Service) GetPoint(ctx context.Context, id uint64) (PickUpPoint, error) {
	return s.repo.Get(ctx, id)
}

// UpdatePoint updates a pick-up point info.
func (s *Service) UpdatePoint(ctx context.Context, point PickUpPoint) error {
	return s.repo.Update(ctx, point)
}

// DeletePoint deletes a pick-up point.
func (s *Service) DeletePoint(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
