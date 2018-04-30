package main

import (
	"github.com/labstack/echo"
	"github.com/bradfitz/gomemcache/memcache"
)

type CustomContext struct {
	echo.Context
	cachedData map[string]*memcache.Item
}

func (c *CustomContext) GetCachedData(key string) *memcache.Item{
	if data, ok := c.cachedData[key]; ok {
		return data
	}

	return nil
}

func (c *CustomContext) SetCachedData(key string, data *memcache.Item) {
	c.cachedData[key] = data
}