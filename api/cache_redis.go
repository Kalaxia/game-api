package api

import(
	"github.com/go-redis/redis/v7"
)

var RedisClient *redis.Client

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "kalaxia_redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}