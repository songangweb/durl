package cache

import (
	"durl/app/share/comm"
	"fmt"
	"github.com/songangweb/mcache"
)

var UrlListCache *urlListCache

type urlListCache struct {
	GoodUrlCache *mcache.ARCCache
	BedUrlCache *mcache.LruCache
}

type Conf struct {
	GoodUrlLen int
	BedUrlLen  int
}

func InitUrlCache(c Conf) {

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

	UrlListCache = &urlListCache{
		GoodUrlCache: goodUrlCache,
		BedUrlCache:  bedUrlCache,
	}
}

func (s *urlListCache) Gget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.GoodUrlCache.Get(key)
	return fullUrl, ok
}

func (s *urlListCache) Gadd(key, value interface{}, expirationTime int64) {
	s.GoodUrlCache.Add(key, value, expirationTime)
}

func (s *urlListCache) Gremove(key interface{}) {
	s.GoodUrlCache.Remove(key)
}

func (s *urlListCache) Bget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.BedUrlCache.Get(key)
	return fullUrl, ok
}

func (s *urlListCache) Badd(key, value interface{}, expirationTime int64) {
	s.BedUrlCache.Add(key, value, expirationTime)
}