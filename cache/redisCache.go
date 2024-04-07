package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/cetnfurkan/core/config"

	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
)

type (
	RedisCache struct {
		cfg    *redisConfig
		client *redis.Client
	}

	redisConfig struct {
		*config.Database
		Extra redisConfigExtra
	}

	redisConfigExtra struct {
		Scheme string `mapstructure:"scheme"`
	}
)

// NewRedisCache creates a new redis cache instance.
//
// It takes a config instance and returns a new cache interface instance.
//
// It will panic
// if it fails to unmarhal extra config data,
// if it fails to parse the connection string or
// if it fails to connect to redis.
func NewRedisCache(cfg *config.Database) Cache {
	redisCache := &RedisCache{
		cfg: &redisConfig{
			Database: cfg,
		},
	}

	redisCache.UnmarshalExtra()

	connectionString := fmt.Sprintf(
		"%s://%s:%s@%s:%d",
		redisCache.cfg.Extra.Scheme,
		redisCache.cfg.User,
		redisCache.cfg.Password,
		redisCache.cfg.Host,
		redisCache.cfg.Port,
	)

	options, err := redis.ParseURL(connectionString)
	if err != nil {
		log.Fatal("Unable to parse redis connection string: ", err)
	}

	redisCache.client = redis.NewClient(options)

	_, err = redisCache.client.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect to redis: ", err)
	}

	return redisCache
}

func (redisCache *RedisCache) UnmarshalExtra() {
	err := mapstructure.Decode(redisCache.cfg.Database.Extra, &redisCache.cfg.Extra)
	if err != nil {
		log.Fatal("Unable to decode extra config into struct: ", err)
	}
}

func (redisCache *RedisCache) Get(key string) (interface{}, error) {
	return redisCache.client.Get(key).Result()
}

func (redisCache *RedisCache) Set(key string, value interface{}, duration time.Duration) error {
	return redisCache.client.Set(key, value, duration).Err()
}

func (redisCache *RedisCache) Client() *redis.Client {
	return redisCache.client
}
