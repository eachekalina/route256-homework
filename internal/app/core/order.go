package core

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"homework/internal/app/order"
	"homework/internal/app/packaging"
)

type OrderCoreService struct {
	orderService      OrderService
	packagingVariants map[packaging.Type]packaging.Packaging
	tracer            trace.Tracer
}

type OrderService interface {
	AddOrder(ctx context.Context, o order.Order) error
	RemoveOrder(ctx context.Context, id uint64) error
	GiveOrders(ctx context.Context, ids []uint64) ([]order.Order, error)
	GetOrders(ctx context.Context, customerId uint64, n int, filterGiven bool) ([]order.Order, error)
	AcceptReturn(ctx context.Context, orderId uint64, customerId uint64) (order.Order, error)
	GetReturns(ctx context.Context, count int, pageNum int) ([]order.Order, error)
}

func NewOrderCoreService(orderService OrderService, packagingTypes map[packaging.Type]packaging.Packaging) *OrderCoreService {
	return &OrderCoreService{
		orderService:      orderService,
		packagingVariants: packagingTypes,
		tracer:            otel.Tracer("internal/app/core/order"),
	}
}
