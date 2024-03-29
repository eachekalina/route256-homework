package core

import (
	"homework/internal/app/order"
	"homework/internal/app/packaging"
)

type OrderCoreService struct {
	orderService      OrderService
	packagingVariants map[packaging.Type]packaging.Variant
}

type OrderService interface {
	AddOrder(o order.Order) error
	RemoveOrder(id uint64) error
	GiveOrders(ids []uint64) error
	GetOrders(customerId uint64, n int, filterGiven bool) ([]order.Order, error)
	AcceptReturn(orderId uint64, customerId uint64) error
	GetReturns(count int, pageNum int) ([]order.Order, error)
}

func NewOrderCoreService(orderService OrderService, packagingVariants map[packaging.Type]packaging.Variant) *OrderCoreService {
	return &OrderCoreService{
		orderService:      orderService,
		packagingVariants: packagingVariants,
	}
}
