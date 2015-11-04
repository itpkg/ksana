package job

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/itpkg/ksana/utils"
)

type RedisStore struct {
	pool *redis.Pool
}

func (p *RedisStore) Push(queue string, msg *Message) error {
	c := p.pool.Get()
	defer c.Close()
	buf, err := utils.ToJson(msg)
	if err != nil {
		return err
	}

	_, err = c.Do("LPUSH", p.queue(queue), buf)
	return err

}

func (p *RedisStore) Pop(queues ...string) (string, *Message, error) {
	c := p.pool.Get()
	defer c.Close()

	args := make([]interface{}, 0)
	for _, q := range queues {
		args = append(args, p.queue(q))
	}
	args = append(args, 0)

	rep, err := redis.Values(c.Do("BRPOP", args...))
	if err != nil {
		return "", nil, err
	}

	name := string(rep[0].([]byte)[6:])
	msg := Message{}
	err = utils.FromJson(rep[1].([]byte), &msg)
	if err != nil {
		return "", nil, err
	}
	return name, &msg, nil

}

func (p *RedisStore) queue(name string) string {
	return fmt.Sprintf("job://%s", name)
}

//==============================================================================
func NewRedisStore(p *redis.Pool) Store {
	return &RedisStore{pool: p}
}
