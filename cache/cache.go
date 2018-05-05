package cache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ohsean53/oceansf/lib"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	EXPIRE = 600
)

type Cache struct {
	*memcache.Client
	cachedData map[string]*Item
}

type Item struct {
	rawCache *memcache.Item
	dirty    bool
}

func NewConnection(endpoint string) *Cache {
	client := memcache.New(endpoint)
	return &Cache{
		client,
		map[string]*Item{},
	}
}

func GetGlobalLockKey(uid int64) string {
	return "globalLock:" + lib.Itoa64(uid)
}

func GetGlobalLockKeyWithId(id string) string {
	return "globalLock:" + id
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

func (c *Cache) Get(key string) ([]byte, error) {
	if localData := c.getCachedData(key); localData != nil {
		log.Debug("get cached from local, key:", key)

		return localData.rawCache.Value, nil
	}

	data, err := c.Client.Get(key)
	if data != nil {
		c.cachedData[key] = &Item{
			rawCache: data,
			dirty:    false,
		}
		return data.Value, err
	}

	return nil, err
}

func (c *Cache) Set(key string, data []byte, expiration int) error {
	err := c.Client.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: int32(expiration),
	}) // embedding 처리 했으므로, 함수이름이 모호하면 이런식으로..호출

	if _, ok := c.cachedData[key]; ok {
		delete(c.cachedData, key)
	}

	return err
}

func (c *Cache) CasDelayed(key string, data []byte, expire int) {
	log.Debug("expire:", expire)
	if _, ok := c.cachedData[key]; !ok {
		c.cachedData[key] = &Item{
			rawCache: &memcache.Item{
				Key: key,
			},
		}
	}

	c.cachedData[key].rawCache.Value = data
	c.cachedData[key].rawCache.Expiration = int32(expire) // 문제 될까?
	c.cachedData[key].dirty = true

	// 로직 수행후 마지막 미들웨어에서 부분에서 commit 처리됨
}

func (c *Cache) CommitAll() {
	for _, item := range c.cachedData {
		if item.dirty == true {
			log.Debug("commit key:", item.rawCache.Key)

			casErr := c.CompareAndSwap(item.rawCache)
			var err error
			if casErr == nil {
				log.Debug("cas success")
				continue
			}

			log.Info("cas err :", casErr)
			if casErr == memcache.ErrCacheMiss {
				// 캐시서버에 없는 경우임
				// c.cachedData[key] 에도 없다고 가정해야함
				// cas 토큰이 없을때는 바로 set 하는게 맞는거같은데 casid가 unexported 되어 있어서 확인불가
				// cas 토큰이 없다하면, cas 기능을 off 했다는 건가? <- 확인 필요.
				// c.cachedData[key] 에 없다면 memcache 에서 get 을 하지 못한 것이므로, cas 를 통해서 저장하지 못함
				// 혹시 모를 캐시 set 경합을 방지하기 위해 add로 처리함 (global lock 유발)
				err = c.Client.Add(item.rawCache)
				log.Info("cas cache miss, add key :", item.rawCache.Key)
				if err == nil {
					log.Debug("add success")
					continue
				}

				log.Info("cas add error :", err)
			}

			// cas 및 add 실패시에 캐시 삭제
			deleteErr := c.Delete(item.rawCache.Key)
			if deleteErr != nil {
				// 삭제도 실패한다면..나중에 오래된 데이터를 불러 올 수 있다.
				// 장애시에 재접속을 하도록 처리하게되면 그만큼 해당 세션의 요청 처리시간이 길어져, 도미노 장애가 생길 수 도 있음
				log.Error("cas fail, delete error:", deleteErr)
			}
		}

	}
}

func (c *Cache) DiscardAll() {
	for key := range c.cachedData {
		delete(c.cachedData, key)
	}
}

func (c *Cache) getCachedData(key string) *Item {
	if data, ok := c.cachedData[key]; ok {
		return data
	}

	return nil
}

func (c *Cache) setCachedData(data *Item) {
	c.cachedData[data.rawCache.Key] = data
}
