package core

import (
	"errors"
	"homework/internal/app/order"
	"homework/internal/app/packaging"
	"time"
)

const dateFormat = "2006-01-02"

type AcceptOrderRequest struct {
	OrderId        uint64
	CustomerId     uint64
	KeepDateString string
	PriceRub       int64
	WeightKg       float64
	PackagingType  string
}

func (s *OrderCoreService) AcceptOrder(req AcceptOrderRequest) error {
	if req.OrderId == 0 {
		return errors.New("valid order id is required")
	}
	if req.CustomerId == 0 {
		return errors.New("valid customer id is required")
	}

	if req.KeepDateString == "" {
		return errors.New("keep date is required")
	}
	keepDate, err := time.ParseInLocation(dateFormat, req.KeepDateString, time.Local)
	if err != nil {
		return err
	}
	keepDate = keepDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	now := time.Now()
	if keepDate.Before(now) {
		return errors.New("keepDate can't be in the past")
	}

	if req.PriceRub <= 0 {
		return errors.New("price must be positive")
	}
	if req.WeightKg <= 0 {
		return errors.New("weight must be positive")
	}

	var packagingVariant packaging.Packaging
	if req.PackagingType != "" {
		var ok bool
		packagingVariant, ok = s.packagingVariants[packaging.Type(req.PackagingType)]
		if !ok {
			return errors.New("invalid packaging type")
		}
	}

	o := order.Order{
		KeepDate:   keepDate,
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
