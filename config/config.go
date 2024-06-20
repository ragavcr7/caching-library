package cache

import (
	"fmt"
	"os"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
type MemcachedConfig struct {
	Addr string
}

func LoadRedisCache() (*RedisConfig, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("address not found")
	}
	password := os.Getenv("REDIS_PASSWORD")
	db := 0 //default db
	dbStr := os.Getenv("REDIS_DB")
	if dbStr != "" {
		_, err := fmt.Sscanf(dbStr, "%d", &db)
		if err != nil {
			fmt.Print("Invalid database value ")
		}
	}
	return &RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
	}, nil
}
func LoadMemCachedCache() (*MemcachedConfig, error) {
	addr := os.Getenv("MEMCACHED_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("address value not set ")
	}
	return &MemcachedConfig{
		Addr: addr,
	}, nil
}
