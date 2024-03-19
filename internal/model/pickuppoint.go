package model

import "fmt"

// PickUpPoint contains fields relevant to a pick-up point
type PickUpPoint struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func (p PickUpPoint) String() string {
	return fmt.Sprintf(
		"%d\t%s\t%s\t%s\n",
		p.Id,
		p.Name,
		p.Address,
		p.Contact)
}
