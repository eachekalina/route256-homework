package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"homework/internal/db"
	"homework/internal/model"
	"homework/internal/storage"
)

// PickUpPointStorage provides a pick-up point storage with a PostgreSQL database as a backend.
type PickUpPointStorage struct {
	db *db.Database
}

// NewPickUpPointStorage returns a new PickUpPointStorage with provided database.
func NewPickUpPointStorage(db *db.Database) *PickUpPointStorage {
	return &PickUpPointStorage{db: db}
}

// Create creates a new pick-up point.
func (s *PickUpPointStorage) Create(ctx context.Context, point model.PickUpPoint) error {
	_, err := s.db.Exec(ctx, "INSERT INTO pickup_points (id, name, address, contact) VALUES ($1, $2, $3, $4);", point.Id, point.Name, point.Address, point.Contact)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.ConstraintName != "" {
		return storage.ErrIdAlreadyExists
	}
	return err
}

// List returns a slice of all pick-up points stored.
func (s *PickUpPointStorage) List(ctx context.Context) ([]model.PickUpPoint, error) {
	var slice []model.PickUpPoint
	err := s.db.Select(ctx, &slice, "SELECT id, name, address, contact FROM pickup_points;")
	return slice, err
}

// Get returns the pick-up point represented by id.
func (s *PickUpPointStorage) Get(ctx context.Context, id uint64) (model.PickUpPoint, error) {
	var point model.PickUpPoint
	err := s.db.Get(ctx, &point, "SELECT id, name, address, contact FROM pickup_points WHERE id = $1;", id)
	if errors.Is(err, sql.ErrNoRows) {
		return point, storage.ErrNoItemFound
	}
	return point, err
}

// Update sets the parameters of a pick-up point to those provided.
func (s *PickUpPointStorage) Update(ctx context.Context, point model.PickUpPoint) error {
	tag, err := s.db.Exec(ctx, "UPDATE pickup_points SET name = $2, address = $3, contact = $4 WHERE id = $1;", point.Id, point.Name, point.Address, point.Contact)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrNoItemFound
	}
	return nil
}

// Delete deletes a pick-up point.
func (s *PickUpPointStorage) Delete(ctx context.Context, id uint64) error {
	tag, err := s.db.Exec(ctx, "DELETE FROM pickup_points WHERE id = $1;", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrNoItemFound
	}
	return nil
}
