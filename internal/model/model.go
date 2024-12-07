package model

import (
	"fmt"
	"github.com/RomanAgaltsev/ya_gophermart/internal/database/queries"
	"net/http"
	"time"
)

type User struct {
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

func (u *User) Bind(r *http.Request) error {
	if u.Login == "" {
		return fmt.Errorf("login is a required field")
	}
	if u.Password == "" {
		return fmt.Errorf("password is a required field")
	}
	return nil
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Order struct {
	Login      string              `db:"login" json:"-"`
	Number     string              `db:"number" json:"number"`
	Status     queries.OrderStatus `db:"status" json:"status"`
	Accrual    float64             `db:"accrual" json:"accrual"`
	UploadedAt time.Time           `db:"uploaded_at" json:"uploaded_at"`
}

type Orders []*Order

type OrderAccrual struct {
	OrderNumber string              `json:"order"`
	Status      queries.OrderStatus `json:"status"`
	Accrual     float64             `json:"accrual"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	OrderNumber string    `db:"order" json:"order"`
	Sum         float64   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at,omitempty"`
}

type Withdrawals []*Withdrawal
