package pickuppoint

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"homework/internal/app/db/mocks"
	"testing"
)

type PostgresRepositoryTestSuite struct {
	suite.Suite
}

func (s *PostgresRepositoryTestSuite) Test_Create() {
	tests := []struct {
		name    string
		point   PickUpPoint
		dbErr   error
		wantErr bool
		err     error
	}{
		{
			name:  "valid",
			point: SamplePickUpPoint,
		},
		{
			name:    "existing id",
			point:   SamplePickUpPoint,
			dbErr:   &pgconn.PgError{ConstraintName: "pickup_points_pkey"},
			wantErr: true,
			err:     ErrIdAlreadyExists,
		},
		{
			name:    "error",
			point:   SamplePickUpPoint,
			dbErr:   assert.AnError,
			wantErr: true,
			err:     assert.AnError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			db := mocks.NewMockDatabase(ctrl)
			repo := NewPostgresRepository(db)
			db.EXPECT().
				Exec(gomock.Any(),
					"INSERT INTO pickup_points (id, name, address, contact) VALUES ($1, $2, $3, $4);",
					tt.point.Id, tt.point.Name, tt.point.Address, tt.point.Contact).
				Return(nil, tt.dbErr)
			err := repo.Create(context.Background(), tt.point)
			if tt.wantErr {
				s.Error(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *PostgresRepositoryTestSuite) Test_List() {
	tests := []struct {
		name    string
		dbErr   error
		want    []PickUpPoint
		wantErr bool
		err     error
	}{
		{
			name: "ok",
			want: SamplePickUpPointSlice,
		},
		{
			name:    "error",
			dbErr:   assert.AnError,
			wantErr: true,
			err:     assert.AnError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			db := mocks.NewMockDatabase(ctrl)
			repo := NewPostgresRepository(db)
			db.EXPECT().
				Select(gomock.Any(), gomock.Any(),
					"SELECT id, name, address, contact FROM pickup_points;").
				DoAndReturn(func(ctx context.Context, dest *[]PickUpPoint, query string, args ...interface{}) error {
					*dest = tt.want
					return tt.err
				})
			point, err := repo.List(context.Background())
			if tt.wantErr {
				s.Error(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)
				s.Equal(tt.want, point)
			}
		})
	}
}

func (s *PostgresRepositoryTestSuite) Test_Get() {
	tests := []struct {
		name    string
		id      uint64
		dbErr   error
		want    PickUpPoint
		wantErr bool
		err     error
	}{
		{
			name: "ok",
			id:   1,
			want: SamplePickUpPoint,
		},
		{
			name:    "not found",
			id:      1,
			dbErr:   pgx.ErrNoRows,
			wantErr: true,
			err:     ErrNoItemFound,
		},
		{
			name:    "error",
			id:      1,
			dbErr:   assert.AnError,
			wantErr: true,
			err:     assert.AnError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			db := mocks.NewMockDatabase(ctrl)
			repo := NewPostgresRepository(db)
			db.EXPECT().
				Get(gomock.Any(), gomock.Any(),
					"SELECT id, name, address, contact FROM pickup_points WHERE id = $1;",
					tt.id).
				DoAndReturn(func(ctx context.Context, dest *PickUpPoint, query string, args ...interface{}) error {
					*dest = tt.want
					return tt.err
				})
			point, err := repo.Get(context.Background(), tt.id)
			if tt.wantErr {
				s.Error(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)
				s.Equal(tt.want, point)
			}
		})
	}
}

func (s *PostgresRepositoryTestSuite) Test_Update() {
	tests := []struct {
		name         string
		point        PickUpPoint
		rowsAffected int64
		dbErr        error
		wantErr      bool
		err          error
	}{
		{
			name:         "ok",
			point:        SamplePickUpPoint,
			rowsAffected: 1,
		},
		{
			name:    "not found",
			point:   SamplePickUpPoint,
			wantErr: true,
			err:     ErrNoItemFound,
		},
		{
			name:    "error",
			point:   SamplePickUpPoint,
			dbErr:   assert.AnError,
			wantErr: true,
			err:     assert.AnError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			db := mocks.NewMockDatabase(ctrl)
			repo := NewPostgresRepository(db)
			tag := mocks.NewMockCommandTag(ctrl)
			tag.EXPECT().RowsAffected().AnyTimes().Return(tt.rowsAffected)
			db.EXPECT().
				Exec(gomock.Any(),
					"UPDATE pickup_points SET name = $2, address = $3, contact = $4 WHERE id = $1;",
					tt.point.Id, tt.point.Name, tt.point.Address, tt.point.Contact).
				Return(tag, tt.dbErr)
			err := repo.Update(context.Background(), tt.point)
			if tt.wantErr {
				s.Error(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *PostgresRepositoryTestSuite) Test_Delete() {
	tests := []struct {
		name         string
		id           uint64
		rowsAffected int64
		dbErr        error
		wantErr      bool
		err          error
	}{
		{
			name:         "ok",
			id:           1,
			rowsAffected: 1,
		},
		{
			name:    "not found",
			id:      1,
			wantErr: true,
			err:     ErrNoItemFound,
		},
		{
			name:    "error",
			id:      1,
			dbErr:   assert.AnError,
			wantErr: true,
			err:     assert.AnError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			db := mocks.NewMockDatabase(ctrl)
			repo := NewPostgresRepository(db)
			tag := mocks.NewMockCommandTag(ctrl)
			tag.EXPECT().RowsAffected().AnyTimes().Return(tt.rowsAffected)
			db.EXPECT().
				Exec(gomock.Any(),
					"DELETE FROM pickup_points WHERE id = $1;",
					tt.id).
				Return(tag, tt.dbErr)
			err := repo.Delete(context.Background(), tt.id)
			if tt.wantErr {
				s.Error(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func TestPostgresRepository(t *testing.T) {
	suite.Run(t, new(PostgresRepositoryTestSuite))
}
