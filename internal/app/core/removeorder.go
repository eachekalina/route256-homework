package core

import (
	"errors"
)

// ReturnOrder removes order associated with provided orderId.
func (s *OrderCoreService) ReturnOrder(orderId uint64) error {
	if orderId == 0 {
		return errors.New("valid order id is required")
	}
	return s.orderService.RemoveOrder(orderId)
}
