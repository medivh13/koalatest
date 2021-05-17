package models

import "time"

type Customers struct {
	CustomerId   string    `db:"customer_id"`
	CustomerName string    `db:"customer_name"`
	Email        string    `db:"email"`
	PhoneNumber  string    `db:"phone_number"`
	Dob          string    `db:"dob"`
	Sex          string    `db:"sex"`
	Salt         string    `db:"salt"`
	Password     string    `db:"password"`
	CreatedDate  time.Time `db:"created_date"`
}

type ExistingCustomer struct {
	CustomerId string `db:"customer_id"`
}

type ExistingCustomerWithSalt struct {
	CustomerId string `db:"customer_id"`
	Salt       string `db:"salt"`
}
