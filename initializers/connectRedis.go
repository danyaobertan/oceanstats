package initializers

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis(config *Config) {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
		//Password: config.RedisPassword,
		DB: config.RedisDb,
	})

	RedisClient = client
}
