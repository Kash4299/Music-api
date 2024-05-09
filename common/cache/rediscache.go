package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		client: client,
	}
}

func valueToString(value any) (string, error) {
	tmp, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(tmp), nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any) error {
	valueStr, err := valueToString(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, valueStr, 10*time.Second).Result()
	return err
}

func (c *RedisCache) SetNoTTL(ctx context.Context, key string, value any) error {
	valueStr, err := valueToString(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, valueStr, 0).Result()
	return err
}

func (c *RedisCache) SetTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	valueStr, err := valueToString(value)
	if err != nil {
		return err
	}
	_, err = c.client.Set(ctx, key, valueStr, ttl).Result()
	return err
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func (c *RedisCache) Close() {
}

func (c *RedisCache) HSet(ctx context.Context, key string, field string, value any) error {
	valueStr, err := valueToString(value)
	if err != nil {
		return err
	}
	data := []any{field, valueStr}
	_, err = c.client.HSet(ctx, key, data).Result()
	return err
}

func (c *RedisCache) HDel(ctx context.Context, key string, field ...string) error {
	_, err := c.client.HDel(ctx, key, field...).Result()
	if err == redis.Nil {
		return nil
	}
	return err
}

func (c *RedisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	value, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	value, err := c.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	_, err := c.client.Del(ctx, key).Result()
	return err
}

func (c *RedisCache) Dels(ctx context.Context, keys []string) error {
	_, err := c.client.Del(ctx, keys...).Result()
	return err
}

func (c *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	slice, err := c.client.Keys(ctx, pattern).Result()
	return slice, err
}

func (c *RedisCache) Incr(ctx context.Context, key string) error {
	_, err := c.client.Incr(ctx, key).Result()
	return err
}

func (c *RedisCache) SetRaw(ctx context.Context, key string, value string) error {
	_, err := c.client.Set(ctx, key, value, redis.KeepTTL).Result()
	return err
}

func (c *RedisCache) HSetRaw(ctx context.Context, key string, field string, value string) error {
	data := []any{field, value}
	_, err := c.client.HSet(ctx, key, data).Result()
	return err
}

func (c *RedisCache) Decr(ctx context.Context, key string) error {
	_, err := c.client.Decr(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("key not found")
	}
	return err
}

func (c *RedisCache) SADD(ctx context.Context, key string, value ...any) error {
	arr := make([]any, 0)
	for _, v := range value {
		str, err := valueToString(v)
		if err != nil {
			return err
		}
		arr = append(arr, str)
	}
	_, err := c.client.SAdd(ctx, key, arr).Result()
	return err
}

func (c *RedisCache) SADDRaw(ctx context.Context, key string, value ...string) error {
	_, err := c.client.SAdd(ctx, key, value).Result()
	return err
}

func (c *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	value, err := c.client.SMembers(ctx, key).Result()
	return value, err
}

func (c *RedisCache) SREM(ctx context.Context, key string, value ...string) error {
	_, err := c.client.SRem(ctx, key, value).Result()
	return err
}

func (c *RedisCache) SPOP(ctx context.Context, key string) (string, error) {
	value, err := c.client.SPop(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

func (c *RedisCache) SPOPN(ctx context.Context, key string, count int64) ([]string, error) {
	value, err := c.client.SPopN(ctx, key, count).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

func (s *RedisCache) SCARD(ctx context.Context, key string) (int64, error) {
	value, err := s.client.SCard(ctx, key).Result()
	return value, err
}
func (s *RedisCache) SRANDMEMBER(ctx context.Context, key string, count int64) ([]string, error) {
	value, err := s.client.SRandMemberN(ctx, key, count).Result()
	return value, err
}

func (c *RedisCache) ZADD(ctx context.Context, key string, score float64, v string) error {
	members := redis.Z{
		Score:  score,
		Member: v,
	}
	_, err := c.client.ZAdd(ctx, key, &members).Result()
	return err
}

func (c *RedisCache) ZRangeByScore(ctx context.Context, key string, min, max float64, count int) ([]string, error) {
	value, err := c.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%f", min),
		Max:    fmt.Sprintf("%f", max),
		Offset: 0,
		Count:  int64(count),
	}).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

func (c *RedisCache) ZRem(ctx context.Context, key string, value ...any) error {
	_, err := c.client.ZRem(ctx, key, value...).Result()
	if err == redis.Nil {
		return nil
	}
	return err
}

func (c *RedisCache) ZCount(ctx context.Context, key string, min, max float64) (int64, error) {
	minStr := fmt.Sprintf("%f", min)
	if min == -1 {
		minStr = "-inf"
	}
	maxStr := fmt.Sprintf("%f", max)
	if max == -1 {
		maxStr = "+inf"
	}
	count, err := c.client.ZCount(ctx, key, minStr, maxStr).Result()
	return count, err
}

func (c *RedisCache) ZScore(ctx context.Context, key string, member string) (float64, error) {
	score, err := c.client.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		return -1, nil
	}
	return score, err
}
