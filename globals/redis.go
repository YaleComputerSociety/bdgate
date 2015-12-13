package globals

import (
	"log"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

const (
	redisDefaultDDb     = 4
	redisDefaultAddress = ":6379"
)

func setupRedis() redis.Conn {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = redisDefaultAddress

	}

	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		dbIndex = redisDefaultDDb
	}

	client, err := redis.Dial("tcp", addr, redis.DialDatabase(dbIndex))
	if err != nil {
		panic("Failed to start redis.\n" + err.Error())
	}

	log.Println("Redis started.")

	return client
}

func closeRedis(conn redis.Conn) {
	conn.Close()
}
