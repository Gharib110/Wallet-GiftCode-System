package models

import "time"

type Transaction struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	// Additional transaction-related fields as needed
}

type TransactionDAO struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	// Additional transaction-related fields as needed
}

type TransactionDTO struct {
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	Amount     float64   `json:"amount"`
	UserCharge float64   `json:"user_charge"`
	Type       string    `json:"type"`
	Timestamp  time.Time `json:"timestamp"`
	// Additional transaction-related fields as needed
}
