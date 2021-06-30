package cache

import "time"

type TTLCache interface {
	Put(key string, value interface{}, duration time.Duration) error
	Get(key string)
	Delete(key string)
}
