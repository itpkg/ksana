package utils

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func OpenRedisPool(host string, port, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if e != nil {
				return nil, e
			}
			_, e = c.Do("SELECT", db)
			if e != nil {
				return nil, e
			}
			return c, e

		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
