package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultReadTimeout  = 6 * time.Second
	defaultWriteTimeout = 6 * time.Second
	defaultDialTimeout  = 10 * time.Second
)

type Redis struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialTimeout  time.Duration

	Client *redis.Client
}

func New(url string, opts ...Option) (*Redis, error) {
	rd := &Redis{
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
		dialTimeout:  defaultDialTimeout,
	}

	for _, opt := range opts {
		opt(rd)
	}

	redisOpts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	redisOpts.ReadTimeout = rd.readTimeout
	redisOpts.WriteTimeout = rd.writeTimeout
	redisOpts.DialTimeout = rd.dialTimeout

	rd.Client = redis.NewClient(redisOpts)

	return rd, nil
}
