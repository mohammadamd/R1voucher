package repositories

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"r1wallet/models"
)

type Voucher interface {
	FindVoucherByCode(code string) (models.VoucherModel, error)
	InsertIntoRedeemedVoucher(userID, voucherID, step int) error
	GetRedeemedCount(voucherID int) (int, error)
	IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error)
	RedeemVoucher(userID int, voucher models.VoucherModel, getStep func(voucher models.VoucherModel) (int, error), success func(userID int, voucher models.VoucherModel) error) error
}

type Redis interface {
	Dequeue(queueName string) (string, error)
	Enqueue(message []byte, queueName string) error
	Increase(key string) (int, error)
	SetValue(key string, value interface{}) error
	GetValue(key string) (string, error)
}

type Repository struct {
	Voucher Voucher
	Redis   Redis
}

func NewRepository(db *sql.DB, re *redis.Client) *Repository {
	return &Repository{
		Voucher: NewVoucherRepository(db),
		Redis:   NewRedisRepository(re),
	}
}
