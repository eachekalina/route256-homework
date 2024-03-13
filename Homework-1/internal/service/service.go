package service

import (
	"Homework-1/internal/model"
	"errors"
	"slices"
	"time"
)

type storage interface {
	Create(order model.Order) error
	List() []model.Order
	Get(id uint64) (model.Order, error)
	Update(order model.Order) error
	Delete(id uint64) error
}

// Service provides methods to work with orders.
type Service struct {
	stor storage
}

// NewService creates a new Service.
func NewService(s storage) Service {
	return Service{stor: s}
}

// AddOrder creates a new order with provided orderId, customerId and keepDate.
func (s *Service) AddOrder(orderId uint64, customerId uint64, keepDate time.Time) error {
	if orderId == 0 {
		return errors.New("valid order id is required")
	}
	if customerId == 0 {
		return errors.New("valid customer id is required")
	}
	now := time.Now()
	if keepDate.Before(now) {
		return errors.New("keepDate can't be in the past")
	}
	err := s.stor.Create(model.Order{
		KeepDate:   keepDate,
		AddDate:    now,
		Id:         orderId,
		CustomerId: customerId,
	})
	if err != nil {
		return err
	}
	return nil
}

// RemoveOrder removes order associated with provided orderId.
func (s *Service) RemoveOrder(orderId uint64) error {
	if orderId == 0 {
		return errors.New("valid order id is required")
	}
	order, err := s.stor.Get(orderId)
	if err != nil {
		return err
	}
	if order.IsGiven && !order.IsReturned {
		return errors.New("order has already been given to customer")
	}
	if order.KeepDate.After(time.Now()) {
		return errors.New("keep date has not arrived yet")
	}
	err = s.stor.Delete(orderId)
	return err
}

// GiveOrders marks orders represented by provided ids as given to customer.
func (s *Service) GiveOrders(orderIds []uint64) error {
	orders := make([]model.Order, len(orderIds))
	now := time.Now()
	var customerId uint64
	for i, id := range orderIds {
		order, err := s.stor.Get(id)
		if err != nil {
			return err
		}
		if order.IsGiven {
			return errors.New("order has already been given")
		}
		if order.KeepDate.Before(now) {
			return errors.New("keep date has already expired")
		}
		if i == 0 {
			customerId = order.CustomerId
		} else if order.CustomerId != customerId {
			return errors.New("orders belong to different customers")
		}
		orders[i] = order
	}
	for _, order := range orders {
		order.IsGiven = true
		order.GiveDate = now
		err := s.stor.Update(order)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetOrders returns slice of orders belonging to customer with provided customerId.
func (s *Service) GetOrders(customerId uint64, n int, filterGiven bool) ([]model.Order, error) {
	if customerId == 0 {
		return nil, errors.New("valid customer id is required")
	}
	if n < 0 {
		return nil, errors.New("n must not be negative")
	}
	l := s.stor.List()
	orders := make([]model.Order, 0)
	for _, order := range l {
		if filterGiven && order.IsGiven && !order.IsReturned {
			continue
		}
		if order.CustomerId != customerId {
			continue
		}
		orders = append(orders, order)
	}
	slices.SortFunc(orders, func(a, b model.Order) int {
		if a.AddDate.Before(b.AddDate) {
			return 1
		}
		if a.AddDate.After(b.AddDate) {
			return -1
		}
		return 0
	})
	if n > 0 && n < len(orders) {
		orders = orders[:n]
	}
	return orders, nil
}

// AcceptReturn marks order as returned by customer.
func (s *Service) AcceptReturn(orderId uint64, customerId uint64) error {
	if orderId == 0 {
		return errors.New("valid order id is required")
	}
	if customerId == 0 {
		return errors.New("valid customer id is required")
	}
	order, err := s.stor.Get(orderId)
	if err != nil {
		return err
	}
	if order.CustomerId != customerId {
		return errors.New("order does not belong to customer")
	}
	if !order.IsGiven {
		return errors.New("order was not given")
	}
	if order.IsReturned {
		return errors.New("order was already returned")
	}
	now := time.Now()
	returnExpirationDate := order.GiveDate.AddDate(0, 0, 2)
	if returnExpirationDate.Before(now) {
		return errors.New("too much time passed since give")
	}
	order.IsReturned = true
	order.ReturnDate = now
	return s.stor.Update(order)
}

// GetReturns returns a slice of orders which were returned by customer.
func (s *Service) GetReturns(count int, pageNum int) ([]model.Order, error) {
	if count <= 0 {
		return nil, errors.New("invalid count of items on page")
	}
	if pageNum < 0 {
		return nil, errors.New("invalid page number")
	}
	l := s.stor.List()
	orders := make([]model.Order, 0)
	for _, order := range l {
		if !order.IsReturned {
			continue
		}
		orders = append(orders, order)
	}
	slices.SortFunc(orders, func(a, b model.Order) int {
		if a.ReturnDate.Before(b.ReturnDate) {
			return 1
		}
		if a.ReturnDate.After(b.ReturnDate) {
			return -1
		}
		return 0
	})
	if len(orders) == 0 && pageNum == 0 {
		return orders, nil
	}
	if pageNum*count >= len(orders) {
		return nil, errors.New("page number is too large")
	}
	if (pageNum+1)*count > len(orders) {
		return orders[pageNum*count:], nil
	}
	return orders[pageNum*count : (pageNum+1)*count], nil
}
