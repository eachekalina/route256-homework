package order

import (
	"errors"
	"slices"
	"time"
)

type Repository interface {
	Create(order Order) error
	List() []Order
	Get(id uint64) (Order, error)
	Update(order Order) error
	Delete(id uint64) error
}

var ErrIdAlreadyExists = errors.New("item with such id already exists")
var ErrNoItemFound = errors.New("no such item found")

type ValidationError struct {
	Err string
}

func (v ValidationError) Error() string {
	return v.Err
}

// Service provides methods to work with orders.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// AddOrder creates a new order with provided orderId, customerId and keepDate.
func (s *Service) AddOrder(o Order) error {
	return s.repo.Create(o)
}

// RemoveOrder removes order associated with provided orderId.
func (s *Service) RemoveOrder(id uint64) error {
	order, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	if order.IsGiven && !order.IsReturned {
		return ValidationError{Err: "order has already been given to customer"}
	}
	if order.KeepDate.After(time.Now()) {
		return ValidationError{Err: "keep date has not arrived yet"}
	}
	err = s.repo.Delete(id)
	return err
}

// GiveOrders marks orders represented by provided ids as given to customer.
func (s *Service) GiveOrders(ids []uint64) error {
	orders := make([]Order, len(ids))
	now := time.Now()
	var customerId uint64
	for i, id := range ids {
		order, err := s.repo.Get(id)
		if err != nil {
			return err
		}
		if order.IsGiven {
			return ValidationError{Err: "order has already been given"}
		}
		if order.KeepDate.Before(now) {
			return ValidationError{Err: "keep date has already expired"}
		}
		if i == 0 {
			customerId = order.CustomerId
		} else if order.CustomerId != customerId {
			return ValidationError{Err: "orders belong to different customers"}
		}
		orders[i] = order
	}
	for _, order := range orders {
		order.IsGiven = true
		order.GiveDate = now
		err := s.repo.Update(order)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetOrders returns slice of orders belonging to customer with provided customerId.
func (s *Service) GetOrders(customerId uint64, n int, filterGiven bool) ([]Order, error) {
	l := s.repo.List()
	orders := make([]Order, 0)
	for _, order := range l {
		if filterGiven && order.IsGiven && !order.IsReturned {
			continue
		}
		if order.CustomerId != customerId {
			continue
		}
		orders = append(orders, order)
	}
	slices.SortFunc(orders, func(a, b Order) int {
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
	order, err := s.repo.Get(orderId)
	if err != nil {
		return err
	}
	if order.CustomerId != customerId {
		return ValidationError{Err: "order does not belong to customer"}
	}
	if !order.IsGiven {
		return ValidationError{Err: "order was not given"}
	}
	if order.IsReturned {
		return ValidationError{Err: "order was already returned"}
	}
	now := time.Now()
	returnExpirationDate := order.GiveDate.AddDate(0, 0, 2)
	if returnExpirationDate.Before(now) {
		return ValidationError{Err: "too much time passed since give"}
	}
	order.IsReturned = true
	order.ReturnDate = now
	return s.repo.Update(order)
}

// ListReturns returns a slice of orders which were returned by customer.
func (s *Service) GetReturns(count int, pageNum int) ([]Order, error) {
	l := s.repo.List()
	orders := make([]Order, 0)
	for _, order := range l {
		if !order.IsReturned {
			continue
		}
		orders = append(orders, order)
	}
	slices.SortFunc(orders, func(a, b Order) int {
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
		return nil, ValidationError{Err: "page number is too large"}
	}
	if (pageNum+1)*count > len(orders) {
		return orders[pageNum*count:], nil
	}
	return orders[pageNum*count : (pageNum+1)*count], nil
}
