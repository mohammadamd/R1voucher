package config

import (
	"github.com/go-redis/redis/v7"
)

func initializeRedis(red Redis) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr: red.Address,
	})

	_, err := r.Ping().Result()
	if err != nil {
		return nil, err
	}

	return r, nil
}
