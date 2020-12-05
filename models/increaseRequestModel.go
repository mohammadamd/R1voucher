package models

type IncreaseRequestModel struct {
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
}
