package pickuppoint

import (
	"context"
	"encoding/json"
	"homework/internal/app/storerr"
	"io"
	"os"
	"sync"
)

const filePerm = 0777

// FileRepository provides a pick-up point Repository with a JSON file as a backend.
type FileRepository struct {
	points   map[uint64]PickUpPoint
	filepath string
	changed  bool
	mutex    sync.RWMutex
}

// NewFileRepository returns a new FileRepository with file stored in the provided path.
func NewFileRepository(path string) (*FileRepository, error) {
	file, err := os.OpenFile(path, os.O_CREATE, filePerm)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	var points map[uint64]PickUpPoint
	if len(bytes) == 0 {
		points = make(map[uint64]PickUpPoint)
	} else {
		err = json.Unmarshal(bytes, &points)
		if err != nil {
			return nil, err
		}
	}
	return &FileRepository{points: points, filepath: path}, nil
}

// Close saves cached pick-up points information into file when needed.
func (s *FileRepository) Close() error {
	s.mutex.RLock()
	if !s.changed {
		s.mutex.RUnlock()
		return nil
	}
	bytes, err := json.Marshal(s.points)
	s.mutex.RUnlock()
	if err != nil {
		return err
	}
	err = os.WriteFile(s.filepath, bytes, filePerm)
	if err != nil {
		return err
	}
	s.changed = false
	return nil
}

// Create creates a new pick-up point.
func (s *FileRepository) Create(ctx context.Context, point PickUpPoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.points[point.Id]
	if exists {
		return storerr.ErrIdAlreadyExists
	}
	s.points[point.Id] = point
	s.changed = true
	return nil
}

// List returns a slice of all pick-up points stored.
func (s *FileRepository) List(ctx context.Context) ([]PickUpPoint, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	slice := make([]PickUpPoint, 0)
	for _, point := range s.points {
		slice = append(slice, point)
	}
	return slice, nil
}

// Get returns the pick-up point represented by id.
func (s *FileRepository) Get(ctx context.Context, id uint64) (PickUpPoint, error) {
	s.mutex.RLock()
	point, found := s.points[id]
	s.mutex.RUnlock()
	if found {
		return point, nil
	}
	return PickUpPoint{}, storerr.ErrNoItemFound
}

// Update sets the parameters of a pick-up point to those provided.
func (s *FileRepository) Update(ctx context.Context, point PickUpPoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, found := s.points[point.Id]
	if !found {
		return storerr.ErrNoItemFound
	}
	s.points[point.Id] = point
	s.changed = true
	return nil
}

// Delete deletes a pick-up point.
func (s *FileRepository) Delete(ctx context.Context, id uint64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, found := s.points[id]
	if !found {
		return storerr.ErrNoItemFound
	}
	delete(s.points, id)
	s.changed = true
	return nil
}
