package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/pastebin/internal/entity"
	"github.com/romankravchuk/pastebin/internal/usecase"
	rds "github.com/romankravchuk/pastebin/pkg/redis"
)

var _ usecase.PastesCache = &PastesCache{}

type PastesCache struct {
	rd *rds.Redis
}

func NewPastesCache(rd *rds.Redis) *PastesCache {
	return &PastesCache{rd: rd}
}

// Create creates paste cache in redis.
// The paste marshalize to slice of bytes.
func (c *PastesCache) Create(ctx context.Context, paste *entity.Paste) error {
	raw, err := paste.MarshalBinary()
	if err != nil {
		return fmt.Errorf("PastesCache.MarshalBinary: %w", err)
	}

	if err := c.rd.Client.Set(ctx, paste.Hash, raw, 0).Err(); err != nil {
		return fmt.Errorf("PastesCache.Redis.Client: %w", err)
	}

	return nil
}

// Delete implements usecase.PastesCache.
func (*PastesCache) Delete(context.Context, string) error {
	panic("unimplemented")
}

// Get returns paste from redis.
// The paste unmarshalize from slice of bytes.
// If paste is not found, returns nil entity and nil error.
func (c *PastesCache) Get(ctx context.Context, hash string) (*entity.Paste, bool, error) {
	cmd := c.rd.Client.Get(ctx, hash)

	raw, err := cmd.Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("PastesCache.Redis.Client: %w", err)
	}

	paste := new(entity.Paste)
	if err := paste.UnmarshalBinary(raw); err != nil {
		return nil, false, fmt.Errorf("PastesCache.UnmarshalBinary: %w", err)
	}

	return paste, true, nil
}
