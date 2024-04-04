package order

import (
	"encoding/json"
	"io"
)

// FileRepository provides an order Repository with a JSON file as a backend.
type FileRepository struct {
	orders  map[uint64]Order
	changed bool
}

// NewFileRepository returns a new FileRepository with file stored in the provided path.
func NewFileRepository(r io.Reader) (*FileRepository, error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var orders map[uint64]Order
	if len(bytes) == 0 {
		orders = make(map[uint64]Order)
	} else {
		err = json.Unmarshal(bytes, &orders)
		if err != nil {
			return nil, err
		}
	}
	return &FileRepository{orders: orders}, nil
}

// Close saves cached orders information into file when needed.
func (s *FileRepository) Close(w io.Writer) error {
	if !s.changed {
		return nil
	}
	bytes, err := json.Marshal(s.orders)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	s.changed = false
	return nil
}

// Create creates a new order.
func (s *FileRepository) Create(order Order) error {
	_, exists := s.orders[order.Id]
	if exists {
		return ErrIdAlreadyExists
	}
	s.orders[order.Id] = order
	s.changed = true
	return nil
}

// List returns a slice of all orders stored.
func (s *FileRepository) List() []Order {
	slice := make([]Order, 0)
	for _, order := range s.orders {
		slice = append(slice, order)
	}
	return slice
}

// Get returns the order represented by id.
func (s *FileRepository) Get(id uint64) (Order, error) {
	if order, found := s.orders[id]; found {
		return order, nil
	}
	return Order{}, ErrNoItemFound
}

// Update sets the parameters of an order to those provided.
func (s *FileRepository) Update(order Order) error {
	_, found := s.orders[order.Id]
	if !found {
		return ErrNoItemFound
	}
	s.orders[order.Id] = order
	s.changed = true
	return nil
}

// Delete deletes an order.
func (s *FileRepository) Delete(id uint64) error {
	_, found := s.orders[id]
	if !found {
		return ErrNoItemFound
	}
	delete(s.orders, id)
	s.changed = true
	return nil
}
