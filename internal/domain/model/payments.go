package model

import "time"

// Payment Банковский платеж
type Payment struct {
	Occ        int       `json:"occ" db:"occ" xml:"occ"`
	Address    string    `json:"address,omitempty" db:"address" xml:"address,omitempty"`
	Date       time.Time `json:"date" db:"date" xml:"date"`
	Value      float64   `json:"value" db:"value" xml:"value"`
	Commission float64   `json:"commission" db:"commission" xml:"commission"`
	Fio        string    `json:"fio,omitempty" db:"fio" xml:"fio,omitempty"`
	// PaymentAccount расчётный счёт
	PaymentAccount string `json:"payment_account,omitempty" db:"payment_account" xml:"payment_account,omitempty"`
}

// Payments структура для хранения платежей
type Payments struct {
	Db []Payment `json:"payments" xml:"payment"`
}
