package model

import "time"

type User struct {
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

type Order struct {
	Number     string    `db:"number" json:"number"`
	Status     string    `db:"status" json:"status"`
	Accrual    float64   `db:"accrual" json:"accrual"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type Orders []*Order

type OrderAccrual struct {
	OrderNumber string  `json:"order"`
	Status      string  ` json:"status"`
	Accrual     float64 `json:"accrual"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
