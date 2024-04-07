package core

import (
	"errors"
	"homework/internal/app/order"
)

type ListReturnsRequest struct {
	Count   int
	PageNum int
}

// ListReturns returns a slice of orders which were returned by customer.
func (s *OrderCoreService) ListReturns(req ListReturnsRequest) ([]order.Order, error) {
	if req.Count <= 0 {
		return nil, errors.New("invalid count of items on page")
	}
	if req.PageNum < 0 {
		return nil, errors.New("invalid page number")
	}
	return s.orderService.GetReturns(req.Count, req.PageNum)
}
