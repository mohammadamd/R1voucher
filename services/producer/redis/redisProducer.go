package redis

import (
	"r1wallet/repositories"
)

type redisProducer struct {
	repository *repositories.Repository
}

func NewRedisProducer(repository *repositories.Repository) *redisProducer {
	return &redisProducer{
		repository: repository,
	}
}

func (r *redisProducer) Produce(message []byte, channelName string) error {
	return r.repository.Redis.Enqueue(message, channelName)
}
