package core

import "context"

// ReturnOrder removes order associated with provided orderId.
func (s *OrderCoreService) ReturnOrder(ctx context.Context, orderId uint64) error {
	ctx, span := s.tracer.Start(ctx, "ReturnOrder")
	defer span.End()

	if orderId == 0 {
		return ValidationError{Err: "valid order id is required"}
	}
	return s.orderService.RemoveOrder(ctx, orderId)
}
