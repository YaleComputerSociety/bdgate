package config

import (
	"log"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type config struct {
	redis *redis.Conn
}

var C *config

func setupRedis() *redis.Conn {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = ":6379"
	}

	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbIndex = 4
	}

	client, err := redis.Dial("tcp", addr, redis.DialDatabase(dbIndex))
	if err != nil {
		panic("Failed to start redis.\n" + err.Error())
	}

	log.Println("Redis started.")

	return &client
}

func Setup() *config {
	C = new(config)

	C.redis = setupRedis()

	return C
}
