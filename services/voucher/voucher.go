package voucher

type Voucher interface {
	Redeem(userID int, code string) error
}
