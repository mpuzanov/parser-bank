package model

import "time"

// Payment Банковский платеж
type Payment struct {
	Occ        int       `json:"occ" db:"occ"`
	Address    string    `json:"address,omitempty" db:"address"`
	Date       time.Time `json:"date" db:"date"`
	Value      float64   `json:"value" db:"value"`
	Commission float64   `json:"commission" db:"commission"`
	Fio        string    `json:"fio,omitempty" db:"fio"`
	// PaymentAccount расчётный счёт
	PaymentAccount string `json:"payment_account,omitempty" db:"payment_account"`
}

// Payments структура для хранения платежей
type Payments struct {
	Db []Payment `json:"payments"`
}
