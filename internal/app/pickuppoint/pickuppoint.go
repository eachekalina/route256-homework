package pickuppoint

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

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

func ListPoints(points []PickUpPoint) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)
	fmt.Fprintf(
		w,
		"%s\t%s\t%s\t%s\n",
		"Id",
		"Name",
		"Address",
		"Contact")
	for _, point := range points {
		fmt.Fprint(w, point)
	}
	w.Flush()
	return buf.String()
}
