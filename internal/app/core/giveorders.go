package core

import (
	"errors"
	"homework/internal/app/order"
)

// GiveOrders marks orders represented by provided ids as given to customer.
func (s *OrderCoreService) GiveOrders(orderIds []uint64) error {
	if orderIds == nil {
		return ValidationError{Err: "list of valid ids is required"}
	}
	err := s.orderService.GiveOrders(orderIds)
	if errors.As(err, &order.ValidationError{}) {
		return ValidationError{Err: err.Error()}
	}
	return err
}
