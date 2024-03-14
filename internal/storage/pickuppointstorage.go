package storage

import (
	"encoding/json"
	"errors"
	"homework/internal/model"
	"io"
	"os"
	"sync"
)

// PickUpPointFileStorage provides a pick-up point storage with a JSON file as a backend.
type PickUpPointFileStorage struct {
	points   map[uint64]model.PickUpPoint
	filepath string
	changed  bool
	mutex    sync.Mutex
}

// NewPickUpPointFileStorage returns a new PickUpPointFileStorage with file stored in the provided path.
func NewPickUpPointFileStorage(path string) (PickUpPointFileStorage, error) {
	file, err := os.OpenFile(path, os.O_CREATE, filePerm)
	if err != nil {
		return PickUpPointFileStorage{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return PickUpPointFileStorage{}, err
	}
	err = file.Close()
	if err != nil {
		return PickUpPointFileStorage{}, err
	}
	var points map[uint64]model.PickUpPoint
	if len(bytes) == 0 {
		points = make(map[uint64]model.PickUpPoint)
	} else {
		err = json.Unmarshal(bytes, &points)
		if err != nil {
			return PickUpPointFileStorage{}, err
		}
	}
	return PickUpPointFileStorage{points: points, filepath: path, mutex: sync.Mutex{}}, nil
}

// Close saves cached pick-up points information into file when needed.
func (s *PickUpPointFileStorage) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.changed {
		return nil
	}
	bytes, err := json.Marshal(s.points)
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
func (s *PickUpPointFileStorage) Create(point model.PickUpPoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.points[point.Id]
	if exists {
		return errors.New("point with such id already exists")
	}
	s.points[point.Id] = point
	s.changed = true
	return nil
}

// List returns a slice of all pick-up points stored.
func (s *PickUpPointFileStorage) List() []model.PickUpPoint {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	slice := make([]model.PickUpPoint, 0)
	for _, point := range s.points {
		slice = append(slice, point)
	}
	return slice
}

// Get returns the pick-up point represented by id.
func (s *PickUpPointFileStorage) Get(id uint64) (model.PickUpPoint, error) {
	s.mutex.Lock()
	point, found := s.points[id]
	s.mutex.Unlock()
	if found {
		return point, nil
	}
	return model.PickUpPoint{}, errors.New("no such point found")
}

// Update sets the parameters of a pick-up point to those provided.
func (s *PickUpPointFileStorage) Update(point model.PickUpPoint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, found := s.points[point.Id]
	if !found {
		return errors.New("no such point found")
	}
	s.points[point.Id] = point
	s.changed = true
	return nil
}

// Delete deletes a pick-up point.
func (s *PickUpPointFileStorage) Delete(id uint64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, found := s.points[id]
	if !found {
		return errors.New("no such point found")
	}
	delete(s.points, id)
	s.changed = true
	return nil
}