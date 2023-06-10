package models2

import (
	"net/http"
	"time"
)

type RedemptionQueueResponse struct {
	Identifier string
	UserRed    *UserRedemption
}

type RedemptionQueueRequest struct {
	R *http.Request
	W http.ResponseWriter
	U *UserInfoDTO
	G *GiftCode
	I string
}

type UserRedemption struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	RedeemedAt time.Time `json:"redeemed_at"`
	Type       string    `json:"type"`
	// Additional user redemption-related fields as needed
}

type RedeemTransactionDTO struct {
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	Amount     float64   `json:"amount"`
	UserCharge float64   `json:"user_charge"`
	Type       string    `json:"type"`
	RedeemedAt time.Time `json:"timestamp"`
	// Additional transaction-related fields as needed
}

type UserRedemptionDAO struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	GiftCodeID int64     `json:"gift_code_id"`
	RedeemedAt time.Time `json:"redeemed_at"`
	Type       string    `json:"type"`
	// Additional user redemption-related fields as needed
}

type UserRedemptionInfo struct {
	PhoneNumber string `json:"phone_number"`
	GiftCode    string `json:"gift_code"`
}

type UserInfoDTO struct {
	ID          float64 `json:"id"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Charge      float64 `json:"charge"`
}
