//go:generate mockgen -source=client.go -mock_names=Client=MockRedisClient -destination=../../mock/mock_redis_client.go -package=mock

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-api/internal/core/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	rc     *clientImpl
	rcOnce sync.Once
)

// RedisClientInterface defines the methods that a Redis client should implement
type Client interface {
	Set(key string, value any, ttl ...int) error
	Get(key string, out any) error
	Del(key string) error
	Close() error
}

type clientImpl struct {
	cfg    *config.Configuration
	client *redis.Client
}

func ProvideClient(cfg *config.Configuration) (Client, error) {
	var err error
	rcOnce.Do(func() {
		addr := fmt.Sprintf("%v:%v", cfg.RedisHost, cfg.RedisPort)
		options := &redis.Options{
			Addr: addr,
			DB:   cfg.RedisDB,
		}
		if len(cfg.RedisPassword) > 0 {
			options.Password = cfg.RedisPassword
		}
		client := redis.NewClient(options)

		var pong string
		pong, err = client.Ping(context.Background()).Result()
		if err != nil {
			logrus.Errorf("error connecting to redis: %v", err)
			return
		} else {
			logrus.Infof("redis connected: %s", pong)
		}

		rc = &clientImpl{
			cfg:    cfg,
			client: client,
		}
	})
	return rc, err
}

func ResetProvideClient() {
	rcOnce = sync.Once{}
}

const defaultTTL = 300 // Adjust this value to set a default TTL

func (r *clientImpl) Set(key string, value any, ttl ...int) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	expiration := defaultTTL
	if len(ttl) > 0 {
		if ttl[0] <= 0 {
			return fmt.Errorf("TTL must be greater than 0")
		}

		expiration = ttl[0]
	}

	timeS := time.Duration(expiration) * time.Second

	var dataByteArray []byte
	var err error

	switch v := value.(type) {
	case []byte:
		dataByteArray = v
	default:
		dataByteArray, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
	}

	return r.client.Set(context.Background(), key, dataByteArray, timeS).Err()
}

func (r *clientImpl) Get(key string, out any) error {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	switch v := out.(type) {
	case *[]byte:
		out = []byte(val)
		return nil
	default:
		return json.Unmarshal([]byte(val), v)
	}
}

func (r *clientImpl) Close() error {
	return r.client.Close()
}

func (r *clientImpl) Del(key string) error {
	return r.client.Del(context.Background(), key).Err()
}
