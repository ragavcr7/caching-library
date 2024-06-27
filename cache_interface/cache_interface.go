// cacheinterface/cache_interface.go

// this file is written independently so it can avoid import cycle issue
package cache_interface

import "time"

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	GetAllKeys() ([]string, error)
	DeleteAllKeys() error
}
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
