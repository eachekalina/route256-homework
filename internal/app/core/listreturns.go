package core

import (
	"context"
	"homework/internal/app/order"
)

type ListReturnsRequest struct {
	Count   int
	PageNum int
}

// ListReturns returns a slice of orders which were returned by customer.
func (s *OrderCoreService) ListReturns(ctx context.Context, req ListReturnsRequest) ([]order.Order, error) {
	ctx, span := s.tracer.Start(ctx, "ListReturns")
	defer span.End()

	if req.Count <= 0 {
		return nil, ValidationError{Err: "invalid count of items on page"}
	}
	if req.PageNum < 0 {
		return nil, ValidationError{Err: "invalid page number"}
	}
	return s.orderService.GetReturns(ctx, req.Count, req.PageNum)
}
