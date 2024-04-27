package core

import (
	"context"
	"errors"
	"homework/internal/app/order"
)

// GiveOrders marks orders represented by provided ids as given to customer.
func (s *OrderCoreService) GiveOrders(ctx context.Context, orderIds []uint64) ([]order.Order, error) {
	ctx, span := s.tracer.Start(ctx, "GiveOrders")
	defer span.End()

	if orderIds == nil {
		return nil, ValidationError{Err: "list of valid ids is required"}
	}
	orders, err := s.orderService.GiveOrders(ctx, orderIds)
	if errors.As(err, &order.ValidationError{}) {
		return nil, ValidationError{Err: err.Error()}
	}
	return orders, err
}
