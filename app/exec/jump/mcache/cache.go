package mcache

import (
	"durl/app/share/comm"
	"fmt"
	"github.com/songangweb/mcache"
)


var NewMcache = new(cache)

type Cache interface {
	Gget(key interface{}) (fullUrl interface{}, ok bool)
	Gadd(key, value interface{}, expirationTime int64)
	Gremove(key, value interface{}, expirationTime int64)

	Bget(key interface{}) (fullUrl interface{}, ok bool)
	Badd(key, value interface{}, expirationTime int64)
}

type cache struct {
	GoodUrlCache *mcache.ARCCache
	BedUrlCache *mcache.LruCache
}

type Conf struct {
	GoodUrlLen int
	BedUrlLen  int
}

func InitCache(c Conf) {

	// 初始化Cache数据池
	goodUrlCache, err := mcache.NewARC(c.GoodUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	// 初始化错误urlCache
	bedUrlCache, err := mcache.NewLRU(c.BedUrlLen)
	if err != nil {
		defer fmt.Println(comm.MsgInitializeCacheError)
		panic(comm.MsgInitializeCacheError + ", err: " + err.Error())
	}

	NewMcache = &cache{
		GoodUrlCache: goodUrlCache,
		BedUrlCache:  bedUrlCache,
	}
}

func (s *cache) Gget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.GoodUrlCache.Get(key)
	return fullUrl, ok
}

func (s *cache) Gadd(key, value interface{}, expirationTime int64) {
	s.GoodUrlCache.Add(key, value, expirationTime)
}

func (s *cache) Gremove(key interface{}) {
	s.GoodUrlCache.Remove(key)
}

func (s *cache) Bget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.BedUrlCache.Get(key)
	return fullUrl, ok
}

func (s *cache) Badd(key, value interface{}, expirationTime int64) {
	s.BedUrlCache.Add(key, value, expirationTime)
}