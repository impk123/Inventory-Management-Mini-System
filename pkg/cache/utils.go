package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// GenerateProductKey สร้าง Key สำหรับเก็บ Cache สินค้า
func GenerateProductKey(id uint) string {
	return fmt.Sprintf("product:%d", id)
}

// SetCached เก็บข้อมูลลง Redis
func SetCached(client *redis.Client, key string, value interface{}, duration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(context.Background(), key, json, duration).Err()
}

// GetCached ดึงข้อมูลจาก Redis
func GetCached(client *redis.Client, key string, dest interface{}) error {
	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
