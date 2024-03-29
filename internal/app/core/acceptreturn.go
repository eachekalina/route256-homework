package core

import (
	"errors"
)

type AcceptReturnRequest struct {
	OrderId    uint64
	CustomerId uint64
}

// AcceptReturn marks order as returned by customer.
func (s *OrderCoreService) AcceptReturn(req AcceptReturnRequest) error {
	if req.OrderId == 0 {
		return errors.New("valid order id is required")
	}
	if req.CustomerId == 0 {
		return errors.New("valid customer id is required")
	}
	return s.orderService.AcceptReturn(req.OrderId, req.CustomerId)
}
