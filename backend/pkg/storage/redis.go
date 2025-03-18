package storage

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(client *redis.Client) *Redis {
	return &Redis{
		client: client,
	}
}

func (r *Redis) GetKey(key string) (*string, error) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	val = strings.Trim(val, "\"")
	return &val, nil
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	val, err := r.client.Keys(pattern).Result()
	if err == redis.Nil || err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, errors.New("not found redis keys by pattern " + pattern)
	}
	return val, nil
}

func (r *Redis) Del(pattern string) (int64, error) {
	val, err := r.client.Del(pattern).Result()
	if err == redis.Nil || err != nil {
		return 0, err
	}
	return val, nil
}

func (r *Redis) SetKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key = strings.Trim(key, "\"")
	err = r.client.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
