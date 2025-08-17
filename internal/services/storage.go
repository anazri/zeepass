package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"

	"github.com/anazri/zeepass/internal/models"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "lu7rodah8aefaiCi",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v. Falling back to in-memory storage.", err)
		log.Println("To use Redis: install Redis server and ensure it's running on localhost:6379")
		rdb = nil
		return
	}
	log.Printf("Connected to Redis: %s", pong)

	// Set Redis client for chat service
	chatService := GetChatService()
	if chatService != nil {
		chatService.SetRedisClient(rdb)
		log.Println("Redis client set for chat service")
	}
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return rdb
}

func StoreEncryptedData(id string, data *models.EncryptedData) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var ttl time.Duration
	if data.ExpiresAt != nil {
		ttl = time.Until(*data.ExpiresAt)
		if ttl <= 0 {
			ttl = time.Minute
		}
	} else {
		ttl = 24 * time.Hour
	}

	key := "zeepass:message:" + id
	err = rdb.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		log.Printf("Redis SET failed for key %s: %v", key, err)
		return err
	}
	log.Printf("Redis SET successful for key %s with TTL %v", key, ttl)
	return nil
}

func GetEncryptedData(id string) (*models.EncryptedData, error) {
	if rdb == nil {
		return nil, fmt.Errorf("Redis not available")
	}

	key := "zeepass:message:" + id
	log.Printf("Redis GET attempt for key: %s", key)
	jsonData, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis GET: key %s not found (redis.Nil)", key)
			return nil, fmt.Errorf("message not found")
		}
		log.Printf("Redis GET error for key %s: %v", key, err)
		return nil, err
	}
	log.Printf("Redis GET successful for key %s, data length: %d", key, len(jsonData))

	var data models.EncryptedData
	err = json.Unmarshal([]byte(jsonData), &data)
	return &data, err
}

func DeleteEncryptedData(id string) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}
	return rdb.Del(ctx, "zeepass:message:"+id).Err()
}

func IncrementViewCount(id string) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}

	data, err := GetEncryptedData(id)
	if err != nil {
		return err
	}

	data.ViewCount++

	if data.ViewCount >= data.MaxViews {
		return DeleteEncryptedData(id)
	}

	return StoreEncryptedData(id, data)
}

func StoreEncryptedFileData(id string, data *models.EncryptedFileData) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var ttl time.Duration
	if data.ExpiresAt != nil {
		ttl = time.Until(*data.ExpiresAt)
		if ttl <= 0 {
			ttl = time.Minute
		}
	} else {
		ttl = 24 * time.Hour
	}

	key := "zeepass:file:" + id
	err = rdb.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		log.Printf("Redis SET failed for file key %s: %v", key, err)
		return err
	}
	log.Printf("Redis SET successful for file key %s with TTL %v", key, ttl)
	return nil
}

func GetEncryptedFileData(id string) (*models.EncryptedFileData, error) {
	if rdb == nil {
		return nil, fmt.Errorf("Redis not available")
	}

	key := "zeepass:file:" + id
	log.Printf("Redis GET attempt for file key: %s", key)
	jsonData, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Redis GET: file key %s not found (redis.Nil)", key)
			return nil, fmt.Errorf("file not found")
		}
		log.Printf("Redis GET error for file key %s: %v", key, err)
		return nil, err
	}
	log.Printf("Redis GET successful for file key %s, data length: %d", key, len(jsonData))

	var data models.EncryptedFileData
	err = json.Unmarshal([]byte(jsonData), &data)
	return &data, err
}

func DeleteEncryptedFileData(id string) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}
	return rdb.Del(ctx, "zeepass:file:"+id).Err()
}

func IncrementFileViewCount(id string) error {
	if rdb == nil {
		return fmt.Errorf("Redis not available")
	}

	data, err := GetEncryptedFileData(id)
	if err != nil {
		return err
	}

	data.ViewCount++

	if data.ViewCount >= data.MaxViews {
		return DeleteEncryptedFileData(id)
	}

	return StoreEncryptedFileData(id, data)
}
