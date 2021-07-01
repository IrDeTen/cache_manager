package cache

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Item struct {
	Value      interface{}
	ValueType  reflect.Type
	CreateTime time.Time
	Expiraton  int64
}

type Cache struct {
	Channel           chan []string
	items             map[string]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func NewCache(defaultExpiration, cleanupInterval time.Duration) *Cache {
	cache := Cache{
		Channel:           make(chan []string),
		items:             make(map[string]Item),
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		go cache.gc()
	}

	return &cache
}

func (c *Cache) Put(key string, value interface{}) error {
	var err error

	if _, found := c.items[key]; found {
		err = errors.New("cache: Item with this key alredy exist")
		return err
	}
	c.items[key] = Item{
		Value:      value,
		ValueType:  reflect.TypeOf(value),
		CreateTime: time.Now(),
		Expiraton:  time.Now().Add(c.defaultExpiration).UnixNano(),
	}
	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	item, found := c.items[key]

	if !found {
		return nil, errors.New("cache: Item with this key does not exist")
	}

	if time.Now().UnixNano() > item.Expiraton {
		return nil, errors.New("cache: Item was expired")

	}

	return item.Value, nil
}

func (c *Cache) GetToObj(key string, dest interface{}) error {
	item, found := c.items[key]
	if !found {
		return errors.New("cache: Item with this key does not exist")
	}

	if time.Now().UnixNano() > item.Expiraton {
		return errors.New("cache: Item was expired")

	}

	if reflect.TypeOf(dest) != item.ValueType {
		err_string := fmt.Sprintf("cache: Destination has a different data type: has %T, expect %T", dest, item.Value)
		return errors.New(err_string)
	}

	dest = item.Value

	return nil
}

func (c *Cache) Delete(key string) error {
	if _, found := c.items[key]; !found {
		return errors.New("cache: Item with this key does not exist")
	}
	delete(c.items, key)
	return nil
}

func (c *Cache) gc() {
	for {
		if len(c.items) > 0 {
			list := make([]string, 0)
			for k, val := range c.items {
				if time.Now().UnixNano() > val.Expiraton {
					delete(c.items, k)
					list = append(list, k)
				}
			}
			c.Channel <- list
		}
		time.Sleep(c.cleanupInterval)
	}
}
