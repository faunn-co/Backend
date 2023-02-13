package orm

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	cache *redis.Client
)

func ConnectRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Error("Error while establishing Live Redis client: %v", err)
	}
	cache = rdb
}

func RedisInstance() *redis.Client {
	if cache == nil {
		ConnectRedis()
	}
	return cache
}

func SET(e echo.Context, key string, data interface{}, ttlSeconds time.Duration) error {
	data, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		log.Errorf("Failed to marshal JSON results: %v\n", jsonErr.Error())
		return jsonErr
	}

	if err := RedisInstance().Set(e.Request().Context(), key, data, ttlSeconds*time.Minute).Err(); err != nil {
		log.Errorf(" Error while writing to redis: %v", err.Error())
		return err
	}
	return nil
}

func GET(e echo.Context, key string) ([]byte, error) {
	val, redisErr := RedisInstance().Get(e.Request().Context(), key).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Warnf("No result of %v in Redis, reading from API", key)
			return nil, nil
		} else {
			log.Errorf("Error while reading from redis: %v", redisErr.Error())
			return nil, redisErr
		}
	}
	log.Infof("Successful | Cached %v", key)
	e.Response().Header().Set("cache", "1")
	return []byte(val), nil
}
