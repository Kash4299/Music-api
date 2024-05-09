package cache

import (
	"context"
	"time"
)

type IMemCache interface {
	Set(key string, value interface{}) error
	SetTTL(key string, value interface{}, t time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Close()
}

type IRedisCache interface {
	Set(ctx context.Context, key string, value interface{}) error
	SetNoTTL(ctx context.Context, key string, value any) error
	SetTTL(ctx context.Context, key string, value interface{}, t time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Dels(ctx context.Context, keys []string) error
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, field ...string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Incr(ctx context.Context, key string) error
	Decr(ctx context.Context, key string) error
	SetRaw(ctx context.Context, key string, value string) error
	HSetRaw(ctx context.Context, key string, field string, value string) error
	Close()
	SCARD(ctx context.Context, key string) (int64, error)
	SPOP(ctx context.Context, key string) (string, error)
	SPOPN(ctx context.Context, key string, count int64) ([]string, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SADD(ctx context.Context, key string, value ...any) error
	SRANDMEMBER(ctx context.Context, key string, count int64) ([]string, error)
	ZADD(ctx context.Context, key string, score float64, v string) error
	ZRangeByScore(ctx context.Context, key string, min, max float64, count int) ([]string, error)
	ZRem(ctx context.Context, key string, value ...any) error
	ZCount(ctx context.Context, key string, min, max float64) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	SADDRaw(ctx context.Context, key string, value ...string) error
	SREM(ctx context.Context, key string, value ...string) error
}

var MCache IMemCache
var RCache IRedisCache
