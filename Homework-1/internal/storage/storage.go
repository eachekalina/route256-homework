package storage

import (
	"Homework-1/internal/model"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type FileStorage struct {
	orders   map[uint64]model.Order
	filepath string
}

func NewFileStorage(path string) (FileStorage, error) {
	file, err := os.OpenFile(path, os.O_CREATE, 0777)
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

func (s *FileStorage) Save() error {
	bytes, err := json.Marshal(s.orders)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.filepath, bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Create(order model.Order) error {
	for _, savedOrder := range s.orders {
		if order.Id == savedOrder.Id {
			return errors.New("order with such id already exists")
		}
	}
	s.orders[order.Id] = order
	return nil
}

func (s *FileStorage) List() []model.Order {
	slice := make([]model.Order, len(s.orders))
	for _, order := range s.orders {
		slice = append(slice, order)
	}
	return slice
}

func (s *FileStorage) Get(id uint64) (model.Order, error) {
	if order, found := s.orders[id]; found {
		return order, nil
	}
	return model.Order{}, errors.New("no such order found")
}

func (s *FileStorage) Update(order model.Order) error {
	for i, savedOrder := range s.orders {
		if order.Id == savedOrder.Id {
			s.orders[i] = order
			return nil
		}
	}
	return errors.New("no such order found")
}

func (s *FileStorage) Delete(id uint64) error {
	for _, savedOrder := range s.orders {
		if id == savedOrder.Id {
			delete(s.orders, savedOrder.Id)
			return nil
		}
	}
	return errors.New("no such order found")
}
