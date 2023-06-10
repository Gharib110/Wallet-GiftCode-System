package models

type User struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Charge      int64  `json:"charge"`
}

type UserDAO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Charge      int64  `json:"charge"`
}
