package models

import (
	"context"
	"github.com/beego/beego/v2/client/cache"
	"time"
)

var (
	urlCache cache.Cache
)

// a wrapper for cache, which can be mocked in tests
type cacheManager interface {
	IsExist(string) bool
	Get(string) interface{}
	Put(string, interface{}, time.Duration) error
}

type caManager struct{}

func (c caManager) IsExist(key string) bool {
	exist, err := urlCache.IsExist(context.TODO(), key)
	if err != nil {
		return false
	}
	return exist
}

func (c caManager) Get(key string) interface{} {
	url, err := urlCache.Get(context.TODO(), key)
	if err != nil {
		return nil
	}
	return url
}

func (c caManager) Put(key string, val interface{}, timeout time.Duration) error {
	return urlCache.Put(context.TODO(), key, val, timeout)
}

var CacheCond cacheManager

func init() {
	if nil == urlCache {
		urlCache, _ = cache.NewCache("memory", `{"interval":3600}`)
	}
	CacheCond = caManager{}
}
