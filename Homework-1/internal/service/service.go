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

type Service struct {
	s storage
}

func New(s storage) Service {
	return Service{s: s}
}

func (s *Service) AddOrder(orderId uint64, customerId uint64, keepDate time.Time) error {
	now := time.Now()
	if keepDate.Before(now) {
		return errors.New("keepDate can't be in the past")
	}
	err := s.s.Create(model.Order{
		GiveDate:   time.Time{},
		ReturnDate: time.Time{},
		KeepDate:   keepDate,
		AddDate:    now,
		Id:         orderId,
		CustomerId: customerId,
		IsGiven:    false,
		IsReturned: false,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveOrder(orderId uint64) error {
	order, err := s.s.Get(orderId)
	if err != nil {
		return err
	}
	if order.IsGiven {
		return errors.New("order has already been given to customer")
	}
	now := time.Now()
	if order.KeepDate.After(now) {
		return errors.New("keep date has not arrived yet")
	}
	err = s.s.Delete(orderId)
	return err
}

func (s *Service) GiveOrders(orderIds []uint64) error {
	orders := make([]model.Order, len(orderIds))
	now := time.Now()
	var customerId uint64
	for i, id := range orderIds {
		order, err := s.s.Get(id)
		if err != nil {
			return nil
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
		err := s.s.Update(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetOrders(customerId uint64, n int, filterGiven bool) ([]model.Order, error) {
	if n < 0 {
		return nil, errors.New("n must not be negative")
	}
	l := s.s.List()
	orders := make([]model.Order, 0)
	for _, order := range l {
		if filterGiven && order.IsGiven {
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
	if n > 0 {
		orders = orders[:n]
	}
	return orders, nil
}

func (s *Service) AcceptReturn(orderId uint64, customerId uint64) error {
	order, err := s.s.Get(orderId)
	if err != nil {
		return err
	}
	if order.CustomerId != customerId {
		return errors.New("order does not belong to customer")
	}
	if !order.IsReturned {
		return errors.New("order was not given")
	}
	now := time.Now()
	returnExpirationDate := order.ReturnDate.AddDate(0, 0, 2)
	if returnExpirationDate.Before(now) {
		return errors.New("too much time passed since give")
	}
	order.IsReturned = true
	order.ReturnDate = now
	err = s.s.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetReturns(count int, pageNum int) ([]model.Order, error) {
	if count <= 0 {
		return nil, errors.New("invalid count of items on page")
	}
	if pageNum < 0 {
		return nil, errors.New("invalid page number")
	}
	l := s.s.List()
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
	return orders[pageNum*count : (pageNum+1)*count], nil
}
