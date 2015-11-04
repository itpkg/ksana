package job

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/ksana/utils"
)

type RedisStore struct {
	pool *redis.Pool
}

func (p *RedisStore) Push(queue string, args ...interface{}) error {
	c := p.pool.Get()
	defer c.Close()
	buf, err := utils.ToJson(args)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", p.queue(queue), buf)
	return err
}

func (p *RedisStore) Pop(queue string, args ...interface{}) error {
	c := p.pool.Get()
	defer c.Close()

	rep, err := redis.Values(c.Do("BRPOP", p.queue(queue), 0))
	if err != nil {
		return err
	}
	return utils.FromJson(rep[1].([]byte), &args)
}

func (p *RedisStore) queue(name string) string {
	return fmt.Sprintf("job://%s", name)
}

//==============================================================================
func NewRedisStore(p *redis.Pool) Store {
	return &RedisStore{pool: p}
}
