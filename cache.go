package cache

import "time"

type TTLCache interface {
	Put(key string, value interface{}, duration time.Duration) error
	Get(key string) (interface{}, error)
	GetToObj(key string, obj *interface{}) error
	Delete(key string) error
	Close()
}
