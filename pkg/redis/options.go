package redis

import "time"

type Option func(*Redis)

func (r *Redis) ReadTimeout(timeout time.Duration) {
	r.readTimeout = timeout
}

func (r *Redis) WriteTimeout(timeout time.Duration) {
	r.writeTimeout = timeout
}

func (r *Redis) DialTimeout(timeout time.Duration) {
	r.dialTimeout = timeout
}
