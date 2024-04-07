package core

import "errors"

// GiveOrders marks orders represented by provided ids as given to customer.
func (s *OrderCoreService) GiveOrders(orderIds []uint64) error {
	if orderIds == nil {
		return errors.New("list of valid ids is required")
	}
	return s.orderService.GiveOrders(orderIds)
}
