package model

import "fmt"

// PickUpPoint contains fields relevant to a pick-up point
type PickUpPoint struct {
	Id      uint64
	Name    string
	Address string
	Contact string
}

func (p PickUpPoint) String() string {
	return fmt.Sprintf(
		"%d\t%s\t%s\t%s\n",
		p.Id,
		p.Name,
		p.Address,
		p.Contact)
}
