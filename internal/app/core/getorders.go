package core

import (
	"context"
	"homework/internal/app/order"
)

type ListOrdersRequest struct {
	CustomerId   uint64
	DisplayCount int
	FilterGiven  bool
}

// ListOrders returns slice of orders belonging to customer with provided customerId.
func (s *OrderCoreService) ListOrders(ctx context.Context, req ListOrdersRequest) ([]order.Order, error) {
	ctx, span := s.tracer.Start(ctx, "ListOrders")
	defer span.End()

	if req.CustomerId == 0 {
		return nil, ValidationError{Err: "valid customer id is required"}
	}
	if req.DisplayCount < 0 {
		return nil, ValidationError{Err: "n must not be negative"}
	}
	return s.orderService.GetOrders(ctx, req.CustomerId, req.DisplayCount, req.FilterGiven)
}
