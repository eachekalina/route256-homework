package model

import "time"

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
