package pickuppoint

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type FileRepositoryTestSuite struct {
	suite.Suite
}

func (s *FileRepositoryTestSuite) Test_Open() {
	tests := []struct {
		name    string
		json    string
		points  map[uint64]PickUpPoint
		wantErr bool
	}{
		{
			name: "valid",
			json: "{\"1\":{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}}",
			points: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
		},
		{
			name:   "empty",
			json:   "",
			points: map[uint64]PickUpPoint{},
		},
		{
			name:    "invalid json",
			json:    "fdsfdsfsd",
			wantErr: true,
		},
		{
			name:    "wrong json",
			json:    "{\"hmm\":\"why\"}",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			r := strings.NewReader(tt.json)
			repo, err := NewFileRepository(r)
			if tt.wantErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
				s.Equal(tt.points, repo.points)
			}
		})
	}
}

func (s *FileRepositoryTestSuite) Test_Close() {
	tests := []struct {
		name    string
		json    string
		points  map[uint64]PickUpPoint
		wantErr bool
	}{
		{
			name: "valid",
			json: "{\"1\":{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}}",
			points: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
		},
		{
			name:   "empty",
			json:   "{}",
			points: map[uint64]PickUpPoint{},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.points
			repo.changed = true
			var buf bytes.Buffer
			err := repo.Close(&buf)
			if tt.wantErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
				s.Equal(tt.json, buf.String())
			}
		})
	}
}

func (s *FileRepositoryTestSuite) Test_Create() {
	tests := []struct {
		name         string
		point        PickUpPoint
		beforePoints map[uint64]PickUpPoint
		afterPoints  map[uint64]PickUpPoint
		changed      bool
		wantErr      bool
		err          error
	}{
		{
			name:         "valid",
			point:        SamplePickUpPoint,
			beforePoints: map[uint64]PickUpPoint{},
			afterPoints: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			changed: true,
		},
		{
			name:  "existing id",
			point: SamplePickUpPoint,
			beforePoints: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			afterPoints: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			changed: false,
			wantErr: true,
			err:     ErrIdAlreadyExists,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.beforePoints
			err := repo.Create(context.Background(), tt.point)
			if tt.wantErr {
				s.NotNil(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.Nil(err)
			}
			s.Equal(tt.afterPoints, repo.points)
		})
	}
}

func (s *FileRepositoryTestSuite) Test_List() {
	tests := []struct {
		name    string
		points  map[uint64]PickUpPoint
		want    []PickUpPoint
		wantErr bool
		err     error
	}{
		{
			name: "ok",
			points: map[uint64]PickUpPoint{
				1: {
					Id:      1,
					Name:    "Generic pick-up point",
					Address: "5, Test st., Moscow",
					Contact: "test@example.com",
				},
				2: {
					Id:      2,
					Name:    "Another pick-up point",
					Address: "19, Sample st., Moscow",
					Contact: "sample@example.com",
				},
			},
			want: SamplePickUpPointSlice,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.points
			points, err := repo.List(context.Background())
			if tt.wantErr {
				s.NotNil(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.Nil(err)
				s.Equal(tt.want, points)
			}
		})
	}
}

func (s *FileRepositoryTestSuite) Test_Get() {
	tests := []struct {
		name    string
		id      uint64
		points  map[uint64]PickUpPoint
		want    PickUpPoint
		wantErr bool
		err     error
	}{
		{
			name: "ok",
			id:   1,
			points: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			want: SamplePickUpPoint,
		},
		{
			name: "not found",
			id:   2,
			points: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			wantErr: true,
			err:     ErrNoItemFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.points
			point, err := repo.Get(context.Background(), tt.id)
			if tt.wantErr {
				s.NotNil(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.Nil(err)
				s.Equal(tt.want, point)
			}
		})
	}
}

func (s *FileRepositoryTestSuite) Test_Update() {
	tests := []struct {
		name         string
		point        PickUpPoint
		beforePoints map[uint64]PickUpPoint
		afterPoints  map[uint64]PickUpPoint
		changed      bool
		wantErr      bool
		err          error
	}{
		{
			name:  "ok",
			point: SamplePickUpPoint,
			beforePoints: map[uint64]PickUpPoint{
				1: {
					Id:      1,
					Name:    "Generic pick-up point before",
					Address: "5, Test st., Moscow before",
					Contact: "test@example.com before",
				},
			},
			afterPoints: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			changed: true,
		},
		{
			name:    "not found",
			point:   SamplePickUpPoint,
			wantErr: true,
			err:     ErrNoItemFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.beforePoints
			err := repo.Update(context.Background(), tt.point)
			if tt.wantErr {
				s.NotNil(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.Nil(err)
			}
			s.Equal(tt.afterPoints, repo.points)
		})
	}
}

func (s *FileRepositoryTestSuite) Test_Delete() {
	tests := []struct {
		name         string
		id           uint64
		beforePoints map[uint64]PickUpPoint
		afterPoints  map[uint64]PickUpPoint
		changed      bool
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			id:   1,
			beforePoints: map[uint64]PickUpPoint{
				1: SamplePickUpPoint,
			},
			afterPoints: map[uint64]PickUpPoint{},
			changed:     true,
		},
		{
			name:    "not found",
			id:      1,
			wantErr: true,
			err:     ErrNoItemFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			repo := &FileRepository{}
			repo.points = tt.beforePoints
			err := repo.Delete(context.Background(), tt.id)
			if tt.wantErr {
				s.NotNil(err)
				s.ErrorIs(err, tt.err)
			} else {
				s.Nil(err)
			}
			s.Equal(tt.afterPoints, repo.points)
		})
	}
}

func TestFileRepository(t *testing.T) {
	suite.Run(t, new(FileRepositoryTestSuite))
}
