package conf

import "github.com/garyburd/redigo/redis"

type config struct {
	Redis redis.Conn
}

var Redis redis.Conn

var C *config

func Setup() *config {
	C = new(config)

	Redis = setupRedis()

	C.Redis = Redis

	return C
}

func Close(c *config) {
	closeRedis(c.Redis)
}
