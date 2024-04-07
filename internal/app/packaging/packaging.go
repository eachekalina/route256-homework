package packaging

import (
	"homework/internal/app/order"
)

type Type string

const (
	BagType  Type = "bag"
	BoxType  Type = "box"
	FilmType Type = "film"
)

type Packaging interface {
	Apply(order order.Order) (order.Order, error)
}
