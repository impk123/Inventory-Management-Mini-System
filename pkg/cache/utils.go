package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ProductVersionKey = "version:product"
	StockVersionKey   = "version:stock"
)

// GenerateProductKey สร้าง Key สำหรับสินค้า (Single)
func GenerateProductKey(id uint) string {
	return fmt.Sprintf("product:id:%d", id)
}

// GenerateStockListKey สร้าง Key สำหรับรายการ Stock (List) โดยอิงตาม Query Params และ Version
func GenerateStockListKey(client *redis.Client, ctx context.Context, params url.Values) string {
	// 1. ดึง Version ปัจจุบัน
	v, err := client.Get(ctx, StockVersionKey).Int()
	if err != nil {
		v = 1
		client.Set(ctx, StockVersionKey, v, 0)
	}

	// 2. เรียง Query Params
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("stock:list:v%d", v))

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(":%s=%s", k, params.Get(k)))
	}

	return sb.String()
}

// GenerateStockHistoryKey สร้าง Key สำหรับประวัติ Stock (History) แยกตาม Product ID
func GenerateStockHistoryKey(client *redis.Client, ctx context.Context, productID int, params url.Values) string {
	// 1. ดึง Version ปัจจุบัน (ใช้ Version เดียวกับ Stock List เพื่อให้ Invalidate พร้อมกัน)
	v, err := client.Get(ctx, StockVersionKey).Int()
	if err != nil {
		v = 1
		client.Set(ctx, StockVersionKey, v, 0)
	}

	// 2. เรียง Query Params
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	// ใช้ Prefix ต่างจาก Stock List เพื่อไม่ให้ Key ชนกัน
	sb.WriteString(fmt.Sprintf("stock:history:v%d:id:%d", v, productID))

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(":%s=%s", k, params.Get(k)))
	}

	return sb.String()
}

// InvalidateStockList "เปลี่ยน Version" เพื่อทำให้ Cache Stock List เดิมทั้งหมดเป็นโมฆะทันที
func InvalidateStockList(client *redis.Client, ctx context.Context) error {
	return client.Incr(ctx, StockVersionKey).Err()
}

// GenerateProductListKey สร้าง Key สำหรับรายการสินค้า (List) โดยอิงตาม Query Params และ Version
func GenerateProductListKey(client *redis.Client, ctx context.Context, params url.Values) string {
	// 1. ดึง Version ปัจจุบัน (ถ้าไม่มีให้สร้างใหม่เริ่มที่ 1)
	v, err := client.Get(ctx, ProductVersionKey).Int()
	if err != nil {
		v = 1
		client.Set(ctx, ProductVersionKey, v, 0)
	}

	// 2. เรียง Query Params เพื่อให้ Key เหมือนเดิมไม่ว่าจะส่ง order ไหนมา
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("product:list:v%d", v)) // แปะ Version

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(":%s=%s", k, params.Get(k)))
	}

	return sb.String()
}

// InvalidateProductList "เปลี่ยน Version" เพื่อทำให้ Cache List เดิมทั้งหมดเป็นโมฆะทันที
func InvalidateProductList(client *redis.Client, ctx context.Context) error {
	return client.Incr(ctx, ProductVersionKey).Err()
}

// SetCached เก็บข้อมูลลง Redis (เพิ่ม ctx)
func SetCached(client *redis.Client, ctx context.Context, key string, value interface{}, duration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, json, duration).Err()
}

// GetCached ดึงข้อมูลจาก Redis (เพิ่ม ctx)
func GetCached(client *redis.Client, ctx context.Context, key string, dest interface{}) error {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
