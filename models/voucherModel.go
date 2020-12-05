package models

type VoucherModel struct {
	ID     int
	Code   string
	Usable int
	Amount int
}
