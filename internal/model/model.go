package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/database/queries"
)

// User is a user structure.
type User struct {
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"password"`
}

// Bind validates user structure.
func (u *User) Bind(r *http.Request) error {
	if u.Login == "" {
		return fmt.Errorf("login is a required field")
	}
	if u.Password == "" {
		return fmt.Errorf("password is a required field")
	}
	return nil
}

// Render tunes rendering of user structure.
func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Order is an order structure.
type Order struct {
	Login      string              `db:"login" json:"-"`
	Number     string              `db:"number" json:"number"`
	Status     queries.OrderStatus `db:"status" json:"status"`
	Accrual    float64             `db:"accrual" json:"accrual"`
	UploadedAt time.Time           `db:"uploaded_at" json:"uploaded_at"`
}

type Orders []*Order

// Render tunes rendering of orders.
func (Orders) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// OrderAccrual is an order accrual structure.
type OrderAccrual struct {
	OrderNumber string              `json:"order"`
	Status      queries.OrderStatus `json:"status"`
	Accrual     float64             `json:"accrual"`
}

// Balance is a user balance structure.
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// Render tunes rendering of balance.
func (*Balance) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Withdrawal is a withdrawal structure.
type Withdrawal struct {
	Login       string    `db:"login" json:"-"`
	OrderNumber string    `db:"order" json:"order"`
	Sum         float64   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at,omitempty"`
}

// Bind validates withdrawal structure.
func (w *Withdrawal) Bind(r *http.Request) error {
	if w.OrderNumber == "" {
		return fmt.Errorf("order is a required field")
	}
	if w.Sum == 0 {
		return fmt.Errorf("sum cannot be equal zero")
	}

	return nil
}

type Withdrawals []*Withdrawal

// Render tunes rendering of withdrawals.
func (Withdrawals) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
