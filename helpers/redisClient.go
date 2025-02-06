package helpers

import (
	"context"
	"encoding/json"
	"event-management/structs"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() *redis.Client {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"), // No password set
		DB:       0,                           // Use default DB
	})

	idk, err := RedisClient.Ping(ctx).Result()
	fmt.Println(idk)
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	return RedisClient
}

func CheckCache(cacheKey string, c *gin.Context) error {
	val, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedPagination structs.Pagination
		json.Unmarshal([]byte(val), &cachedPagination)
		c.JSON(http.StatusOK, gin.H{
			"result": cachedPagination,
		})
	}
	return err
}

func SetCache(pagination interface{}, cacheKey string) {
	jsonData, _ := json.Marshal(pagination)
	RedisClient.Set(ctx, cacheKey, jsonData, time.Minute*5)
}

func InvalidateCache(tableName string) {
	keys, err := RedisClient.Keys(ctx, tableName+":page:*").Result()
	if err != nil {
		return
	}

	if len(keys) > 0 {
		RedisClient.Del(ctx, keys...)
	}
}
