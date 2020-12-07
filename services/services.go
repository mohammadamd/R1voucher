package services

import (
	"r1wallet/config"
	"r1wallet/repositories"
	"r1wallet/services/producer"
	"r1wallet/services/producer/redis"
	"r1wallet/services/voucher"
	"r1wallet/services/voucher/creditVoucher"
)

type Services struct {
	Producer producer.Producer
	Voucher  voucher.Voucher
}

func NewServices(repository *repositories.Repository, app *config.ConfiguredApp) *Services {
	return &Services{
		Voucher:  creditVoucher.NewCreditVoucher(repository, app.Config.App.ComQueueName),
		Producer: redis.NewRedisProducer(repository),
	}
}
