package models

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

// open redis connection
func RedisInit() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("redis_host"), os.Getenv("redis_port")),
	})
	pong, err := Rdb.Ping(context.TODO()).Result()
	fmt.Println("Redis is responding :: ", pong, err)
}
