package main

import (
	"time"

	"github.com/IrDeTen/cache_manager"
)

func main() {
	type Test struct {
		Name string
		Age  uint
	}

	user := Test{
		Name: "test_user",
		Age:  20,
	}
	cache := cache_manager.NewCache(15*time.Minute, 15*time.Second)

	cache.Put("test", user)
	outUser2 := new(Test)
	cache.GetToObj("test", outUser2)
}
