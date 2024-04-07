package pickuppoint

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"homework/internal/app/db"
)

// PostgresRepository provides a pick-up point Repository with a PostgreSQL database as a backend.
type PostgresRepository struct {
	db *db.Database
}

// NewPostgresRepository returns a new PostgresRepository with provided database.
func NewPostgresRepository(db *db.Database) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create creates a new pick-up point.
func (s *PostgresRepository) Create(ctx context.Context, point PickUpPoint) error {
	_, err := s.db.Exec(ctx, "INSERT INTO pickup_points (id, name, address, contact) VALUES ($1, $2, $3, $4);", point.Id, point.Name, point.Address, point.Contact)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.ConstraintName != "" {
		return ErrIdAlreadyExists
	}
	return err
}

// List returns a slice of all pick-up points stored.
func (s *PostgresRepository) List(ctx context.Context) ([]PickUpPoint, error) {
	var slice []PickUpPoint
	err := s.db.Select(ctx, &slice, "SELECT id, name, address, contact FROM pickup_points;")
	if err != nil {
		return nil, err
	}
	return slice, nil
}

// Get returns the pick-up point represented by id.
func (s *PostgresRepository) Get(ctx context.Context, id uint64) (PickUpPoint, error) {
	var point PickUpPoint
	err := s.db.Get(ctx, &point, "SELECT id, name, address, contact FROM pickup_points WHERE id = $1;", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return point, ErrNoItemFound
	}
	return point, err
}

// Update sets the parameters of a pick-up point to those provided.
func (s *PostgresRepository) Update(ctx context.Context, point PickUpPoint) error {
	tag, err := s.db.Exec(ctx, "UPDATE pickup_points SET name = $2, address = $3, contact = $4 WHERE id = $1;", point.Id, point.Name, point.Address, point.Contact)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoItemFound
	}
	return nil
}

// Delete deletes a pick-up point.
func (s *PostgresRepository) Delete(ctx context.Context, id uint64) error {
	tag, err := s.db.Exec(ctx, "DELETE FROM pickup_points WHERE id = $1;", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoItemFound
	}
	return nil
}
