package repositories

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"r1wallet/models"
)

type Voucher interface {
	FindVoucherByCode(code string) (models.VoucherModel, error)
	InsertIntoRedeemedVoucher(userID, voucherID int) error
	GetRedeemedCount(code string) (int, error)
	IsUserRedeemedVoucherBefore(userID, voucherID int) (bool, error)
	RedeemVoucher(userID int, voucher models.VoucherModel, closure func(userID int, voucher models.VoucherModel) error) error
}

type Redis interface {
	SubscribeChannel(channelName string) <-chan *redis.Message
	PublishInChannel(message []byte, channelName string) error
	Increase(key string) error
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
