package data_access

import (
	"../config"
	"fmt"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnectToRedis(host string, port int, password string, db int) {

	addr := fmt.Sprintf("%s:%d", host, port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func init() {
	ConnectToRedis(
		config.CONFIG.Redis.Host,
		config.CONFIG.Redis.Port,
		config.CONFIG.Redis.Password,
		config.CONFIG.Redis.DB,
	)
}
