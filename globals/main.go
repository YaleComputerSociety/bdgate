package globals

import "github.com/garyburd/redigo/redis"

type globals struct {
	Redis redis.Conn
}

var Redis redis.Conn

var C *globals

func Setup() *globals {
	C = new(globals)
	Redis = setupRedis()
	C.Redis = Redis

	return C
}

func Close(c *globals) {
	closeRedis(c.Redis)
}
