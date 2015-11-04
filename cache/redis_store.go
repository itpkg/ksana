package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisStore struct {
	Pool *redis.Pool `inject:""`
}

func (p *RedisStore) Get(key string, val interface{}) error {
	c := p.Pool.Get()
	defer c.Close()

	buf, err := redis.Bytes(c.Do("GET", p.key(key)))
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, val)
}

func (p *RedisStore) Set(key string, val interface{}, exp time.Duration) error {
	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c := p.Pool.Get()
	defer c.Close()

	key = p.key(key)
	if exp > 0 {
		_, err = c.Do("SET", key, buf, "EX", int(exp.Seconds()))
	} else {
		_, err = c.Do("SET", key, buf)
	}
	return err
}

func (p *RedisStore) Delete(key string) error {
	c := p.Pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", p.key(key))
	return err
}

func (p *RedisStore) Flush() error {
	c := p.Pool.Get()
	defer c.Close()
	keys, err := redis.Values(c.Do("KEYS", p.key("*")))
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = c.Do("DEL", keys...)
	}
	return err
}

func (p *RedisStore) key(k string) string {
	return fmt.Sprintf("cache://%s/", k)
}
