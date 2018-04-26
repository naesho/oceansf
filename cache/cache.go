package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
)

var (
	Mcache *memcache.Client
)

func InitMemcache(endpoint string) {
	Mcache = memcache.New(endpoint)
	log.Debug("memcached init")
}