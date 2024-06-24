// cacheinterface/cache_interface.go
package cacheinterface

import "time"

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	GetAllKeys() ([]string, error)
	DeleteAllKeys() error
}
