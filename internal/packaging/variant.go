package packaging

import "homework/internal/model"

type Type string

const (
	BagType  Type = "bag"
	BoxType  Type = "box"
	WrapType Type = "wrap"
)

type Variant interface {
	Apply(order model.Order) (model.Order, error)
}
