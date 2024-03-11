package storage

import (
	"encoding/json"
	"errors"
	"homework/internal/model"
	"io"
	"os"
)

// FileStorage provides a storage with a JSON file as a backend.
type FileStorage struct {
	orders   map[uint64]model.Order
	filepath string
	changed  bool
}

const filePerm = 0777

// NewFileStorage returns a new FileStorage with file stored in the provided path.
func NewFileStorage(path string) (FileStorage, error) {
	file, err := os.OpenFile(path, os.O_CREATE, filePerm)
	if err != nil {
		return FileStorage{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return FileStorage{}, err
	}
	err = file.Close()
	if err != nil {
		return FileStorage{}, err
	}
	var orders map[uint64]model.Order
	if len(bytes) == 0 {
		orders = make(map[uint64]model.Order)
	} else {
		err = json.Unmarshal(bytes, &orders)
		if err != nil {
			return FileStorage{}, err
		}
	}
	return FileStorage{orders: orders, filepath: path}, nil
}

// Close saves cached order information into file when needed.
func (s *FileStorage) Close() error {
	if !s.changed {
		return nil
	}
	bytes, err := json.Marshal(s.orders)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.filepath, bytes, filePerm)
	if err != nil {
		return err
	}
	return nil
}

// Create creates a new order.
func (s *FileStorage) Create(order model.Order) error {
	_, exists := s.orders[order.Id]
	if exists {
		return errors.New("order with such id already exists")
	}
	s.orders[order.Id] = order
	s.changed = true
	return nil
}

// List returns a slice of all orders stored.
func (s *FileStorage) List() []model.Order {
	slice := make([]model.Order, len(s.orders))
	for _, order := range s.orders {
		slice = append(slice, order)
	}
	return slice
}

// Get returns the order represented by id.
func (s *FileStorage) Get(id uint64) (model.Order, error) {
	if order, found := s.orders[id]; found {
		return order, nil
	}
	return model.Order{}, errors.New("no such order found")
}

// Update sets the parameters of an order to those provided.
func (s *FileStorage) Update(order model.Order) error {
	_, found := s.orders[order.Id]
	if !found {
		return errors.New("no such order found")
	}
	s.orders[order.Id] = order
	return nil
}

// Delete deletes an order.
func (s *FileStorage) Delete(id uint64) error {
	for _, savedOrder := range s.orders {
		if id == savedOrder.Id {
			delete(s.orders, savedOrder.Id)
			s.changed = true
			return nil
		}
	}
	return errors.New("no such order found")
}
