package models

type RedeemVoucherRequest struct {
	UserID int    `json:"user_id"`
	Code   string `json:"code"`
}
