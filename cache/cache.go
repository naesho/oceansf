package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	"github.com/naesho/oceansf/lib"
)

var (
	Mcache *memcache.Client
)

type Cache struct {
	*memcache.Client
	// delay 작업 큐..?
}

func InitMemcache(endpoint string) {
	Mcache = memcache.New(endpoint)
	log.Debug("memcached init")
}

func GetGlobalLockKey(uid int64) string {
	return "globalLock:" + lib.Itoa64(uid)
}

func (c *Cache) Lock(key string) {
	c.Add(&memcache.Item{
		Key:key,
		Value:[]byte("1"),
		Expiration: 10,
	})
}

func (c *Cache) UnLock(key string) {
	c.Delete(key)
}