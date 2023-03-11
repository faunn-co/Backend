package orm

import (
	"context"
	"encoding/json"
	"github.com/aaronangxz/AffiliateManager/logger"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var (
	cache      *redis.Client
	REDIS_HOST string
	REDIS_PASS string
)

func ConnectRedis() {
	err := godotenv.Load(getEnvDir())
	if err != nil {
		logger.Warn(context.Background(), "Error loading .env file")
	}

	ENV = os.Getenv("ENV")
	switch ENV {
	case "PROD":
		REDIS_HOST = os.Getenv("PROD_REDIS_HOST")
		REDIS_PASS = os.Getenv("PROD_REDIS_PASS")
		logger.Info(context.Background(), "Connecting to PROD Redis")
		break
	default:
		REDIS_HOST = os.Getenv("TEST_REDIS_HOST")
		REDIS_PASS = os.Getenv("TEST_REDIS_PASS")
		logger.Info(context.Background(), "Connecting to TEST Redis")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST,
		Password: REDIS_PASS,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.ErrorMsg(context.Background(), "Error while establishing Redis client: %v", err.Error())
	} else {
		logger.Info(context.Background(), "Successfully connected to redis")
	}
	cache = rdb
}

func RedisInstance() *redis.Client {
	if cache == nil {
		ConnectRedis()
	}
	return cache
}

func SET(c context.Context, key string, data interface{}, ttlSeconds time.Duration) error {
	data, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		logger.ErrorMsg(c, "Failed to marshal JSON results: %v\n", jsonErr.Error())
		return jsonErr
	}

	if err := RedisInstance().Set(c, key, data, ttlSeconds*time.Minute).Err(); err != nil {
		logger.ErrorMsg(c, " Error while writing to redis: %v", err.Error())
		return err
	}
	logger.Info(c, "Successfully written to redis: %v", key)
	return nil
}

func GET(e echo.Context, c context.Context, key string, needResponseHeader bool) ([]byte, error) {
	val, redisErr := RedisInstance().Get(c, key).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			logger.Warn(c, "No result of %v in Redis, reading from API", key)
			return nil, nil
		} else {
			logger.ErrorMsg(c, "Error while reading from redis: %v", redisErr.Error())
			return nil, redisErr
		}
	}
	logger.Info(c, "Successful | Cached %v", key)
	if needResponseHeader {
		e.Response().Header().Set("X-From-Cache", "1")
	}
	return []byte(val), nil
}
