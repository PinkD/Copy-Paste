package cpst

import (
	"github.com/go-redis/redis"
	"runtime"
	"strconv"
	"sync/atomic"
)

type redisClient struct {
	r *redis.Client
	c uint64
}

func newRedis(addr string) *redisClient {
	return &redisClient{
		r: redis.NewClient(&redis.Options{
			Addr:     addr,
			PoolSize: runtime.NumCPU(),
		}),
		c: 0,
	}
}

func (r *redisClient) genCode() (uint64, error) {
	count := atomic.AddUint64(&r.c, 1) - 1
	return count, nil
}

func (r *redisClient) ContainsContent(sha, content string) (code uint64, err error) {
	codeStr, err := r.r.Get(sha).Result()
	if err != nil || len(codeStr) == 0 {
		return //no sha
	}
	code, _ = strconv.ParseUint(codeStr, 10, 64)
	_content, err := r.r.Get(codeStr).Result()
	if len(_content) == len(content) && _content == content {
		return //same content
	}
	return 0, err //sha collision
}

func (r *redisClient) SaveContent(data *contentData) error {
	r.r.Set(string(data.Code), data.Content, 0)
	r.r.Set(data.Sha, data.Code, 0)
	return nil
}

func (r *redisClient) GetContent(code uint64) (string, error) {
	return r.r.Get(string(code)).Result()
}

func (r *redisClient) setCount(count uint64) {
	r.c = count
}
