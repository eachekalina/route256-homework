package pickuppoint

import "fmt"

// PickUpPoint contains fields relevant to a pick-up point
type PickUpPoint struct {
	Id      uint64 `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	Contact string `json:"contact" db:"contact"`
}

func (p PickUpPoint) String() string {
	return fmt.Sprintf(
		"%d\t%s\t%s\t%s\n",
		p.Id,
		p.Name,
		p.Address,
		p.Contact)
}
