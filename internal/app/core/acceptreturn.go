package core

import (
	"context"
	"homework/internal/app/order"
)

type AcceptReturnRequest struct {
	OrderId    uint64
	CustomerId uint64
}

// AcceptReturn marks order as returned by customer.
func (s *OrderCoreService) AcceptReturn(ctx context.Context, req AcceptReturnRequest) (order.Order, error) {
	ctx, span := s.tracer.Start(ctx, "AcceptReturn")
	defer span.End()

	if req.OrderId == 0 {
		return order.Order{}, ValidationError{Err: "valid order id is required"}
	}
	if req.CustomerId == 0 {
		return order.Order{}, ValidationError{Err: "valid customer id is required"}
	}
	return s.orderService.AcceptReturn(ctx, req.OrderId, req.CustomerId)
}
