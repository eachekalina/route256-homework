package file

import (
	"encoding/json"
	"homework/internal/model"
	"homework/internal/storage"
	"io"
	"os"
)

// OrderFileStorage provides an order storage with a JSON file as a backend.
type OrderFileStorage struct {
	orders   map[uint64]model.Order
	filepath string
	changed  bool
}

// NewOrderFileStorage returns a new OrderFileStorage with file stored in the provided path.
func NewOrderFileStorage(path string) (OrderFileStorage, error) {
	file, err := os.OpenFile(path, os.O_CREATE, filePerm)
	if err != nil {
		return OrderFileStorage{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return OrderFileStorage{}, err
	}
	err = file.Close()
	if err != nil {
		return OrderFileStorage{}, err
	}
	var orders map[uint64]model.Order
	if len(bytes) == 0 {
		orders = make(map[uint64]model.Order)
	} else {
		err = json.Unmarshal(bytes, &orders)
		if err != nil {
			return OrderFileStorage{}, err
		}
	}
	return OrderFileStorage{orders: orders, filepath: path}, nil
}

// Close saves cached orders information into file when needed.
func (s *OrderFileStorage) Close() error {
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
	s.changed = false
	return nil
}

// Create creates a new order.
func (s *OrderFileStorage) Create(order model.Order) error {
	_, exists := s.orders[order.Id]
	if exists {
		return storage.ErrIdAlreadyExists
	}
	s.orders[order.Id] = order
	s.changed = true
	return nil
}

// List returns a slice of all orders stored.
func (s *OrderFileStorage) List() []model.Order {
	slice := make([]model.Order, 0)
	for _, order := range s.orders {
		slice = append(slice, order)
	}
	return slice
}

// Get returns the order represented by id.
func (s *OrderFileStorage) Get(id uint64) (model.Order, error) {
	if order, found := s.orders[id]; found {
		return order, nil
	}
	return model.Order{}, storage.ErrNoItemFound
}

// Update sets the parameters of an order to those provided.
func (s *OrderFileStorage) Update(order model.Order) error {
	_, found := s.orders[order.Id]
	if !found {
		return storage.ErrNoItemFound
	}
	s.orders[order.Id] = order
	s.changed = true
	return nil
}

// Delete deletes an order.
func (s *OrderFileStorage) Delete(id uint64) error {
	_, found := s.orders[id]
	if !found {
		return storage.ErrNoItemFound
	}
	delete(s.orders, id)
	s.changed = true
	return nil
}
