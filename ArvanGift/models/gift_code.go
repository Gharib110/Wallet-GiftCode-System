package models2

import "time"

type GiftCode struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Amount          float64 `json:"amount"`
	IsActive        bool    `json:"is_active"`
	RedemptionLimit int     `json:"redemption_limit"`
	RedemptionCount int     `json:"redemption_count"`
	StartTime       string  `json:"start_time"`
	ExpirationTime  string  `json:"expiration"`
}

type GiftCodeDAO struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Amount          float64   `json:"amount"`
	IsActive        bool      `json:"is_active"`
	RedemptionLimit int       `json:"redemption_limit"`
	RedemptionCount int       `json:"redemption_count"`
	StartTime       time.Time `json:"start_time"`
	ExpirationTime  time.Time `json:"expiration"`
}
