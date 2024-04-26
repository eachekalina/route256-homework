package core

import (
	"homework/internal/app/order"
	"homework/internal/app/packaging"
	"time"
)

type AcceptOrderRequest struct {
	OrderId       uint64
	CustomerId    uint64
	KeepDate      time.Time
	PriceRub      int64
	WeightKg      float64
	PackagingType string
}

func (s *OrderCoreService) AcceptOrder(req AcceptOrderRequest) error {
	if req.OrderId == 0 {
		return ValidationError{Err: "valid order id is required"}
	}
	if req.CustomerId == 0 {
		return ValidationError{Err: "valid customer id is required"}
	}

	if req.KeepDate.IsZero() {
		return ValidationError{Err: "keep date is required"}
	}
	now := time.Now()
	if req.KeepDate.Before(now) {
		return ValidationError{Err: "keepDate can't be in the past"}
	}

	if req.PriceRub <= 0 {
		return ValidationError{Err: "price must be positive"}
	}
	if req.WeightKg <= 0 {
		return ValidationError{Err: "weight must be positive"}
	}

	var packagingVariant packaging.Packaging
	if req.PackagingType != "" {
		var ok bool
		packagingVariant, ok = s.packagingVariants[packaging.Type(req.PackagingType)]
		if !ok {
			return ValidationError{Err: "invalid packaging type"}
		}
	}

	o := order.Order{
		KeepDate:   req.KeepDate,
		AddDate:    now,
		Id:         req.OrderId,
		CustomerId: req.CustomerId,
		PriceRub:   req.PriceRub,
		WeightKg:   req.WeightKg,
	}
	if packagingVariant != nil {
		var err error
		o, err = packagingVariant.Apply(o)
		if err != nil {
			return err
		}
	}
	return s.orderService.AddOrder(o)
}
