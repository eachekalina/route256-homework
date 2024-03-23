package model

import (
	"fmt"
	"time"
)

// Order contains fields relevant to an order.
type Order struct {
	GiveDate   time.Time `json:"give_date"`
	ReturnDate time.Time `json:"return_date"`
	KeepDate   time.Time `json:"keep_date"`
	AddDate    time.Time `json:"add_date"`
	Id         uint64    `json:"id"`
	CustomerId uint64    `json:"customer_id"`
	IsGiven    bool      `json:"is_given"`
	IsReturned bool      `json:"is_returned"`
}

const dateFormat = "2006-01-02"

func (o Order) String() string {
	displayedGiveDate := "-"
	if o.IsGiven {
		displayedGiveDate = o.GiveDate.Format(dateFormat)
	}
	displayedReturnDate := "-"
	if o.IsReturned {
		displayedReturnDate = o.ReturnDate.Format(dateFormat)
	}
	return fmt.Sprintf(
		"%d\t%d\t%s\t%s\t%t\t%s\t%t\t%s\n",
		o.Id,
		o.CustomerId,
		o.AddDate.Format(dateFormat),
		o.KeepDate.Format(dateFormat),
		o.IsGiven,
		displayedGiveDate,
		o.IsReturned,
		displayedReturnDate)
}
