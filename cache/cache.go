package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
	"github.com/naesho/oceansf/lib"
	"time"
	"errors"
)

var (
	Mcache *memcache.Client
)

type Cache struct {
	*memcache.Client
	cachedData map[string]*memcache.Item
}

func InitMemcache(endpoint string) {
	Mcache = memcache.New(endpoint)
	log.Debug("memcached init")
}

func NewConnection(endpoint string) *Cache {
	return &Cache{
		memcache.New(endpoint),
		map[string]*memcache.Item{},
	}
}

func GetGlobalLockKey(uid int64) string {
	return "globalLock:" + lib.Itoa64(uid)
}

func (c *Cache) Lock(key string) error {

	maxCount := 20 // TODO config 로 나중에..

	for i := 0; i < maxCount; i++ {
		err := c.Add(&memcache.Item{
			Key:        key,
			Value:      []byte("1"),
			Expiration: 10,
		})

		if err == memcache.ErrNotStored {
			// spin lock
			log.Debug("lock key:", key, "try :", i)
			time.Sleep(time.Millisecond * 20) // TODO spin lock 시간 config
			continue
		} else {
			log.Error()
		}

		return nil
	}

	return errors.New("lock fail")
}

func (c *Cache) UnLock(key string) {
	c.Delete(key)
}

// cas delayed
// commit all

func (c *Cache) GetCachedData(key string) *memcache.Item{
	if data, ok := c.cachedData[key]; ok {
		return data
	}

	return nil
}

func (c *Cache) SetCachedData(key string, data *memcache.Item) {
	c.cachedData[key] = data
}