package core

import (
	"homework/internal/app/order"
)

type ListOrdersRequest struct {
	CustomerId   uint64
	DisplayCount int
	FilterGiven  bool
}

// ListOrders returns slice of orders belonging to customer with provided customerId.
func (s *OrderCoreService) ListOrders(req ListOrdersRequest) ([]order.Order, error) {
	if req.CustomerId == 0 {
		return nil, ValidationError{Err: "valid customer id is required"}
	}
	if req.DisplayCount < 0 {
		return nil, ValidationError{Err: "n must not be negative"}
	}
	return s.orderService.GetOrders(req.CustomerId, req.DisplayCount, req.FilterGiven)
}
