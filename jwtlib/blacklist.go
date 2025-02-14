package jwtlib

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

func BlacklistToken(token string, expiry time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, "blacklist:"+token, "1", expiry).Err()
}

func IsTokenBlacklisted(token string) bool {
	ctx := context.Background()
	_, err := rdb.Get(ctx, "blacklist:"+token).Result()
	return err == nil
}
