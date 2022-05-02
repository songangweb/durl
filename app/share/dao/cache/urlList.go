package cache

import (
	"fmt"

	"durl/app/share/comm"

	"github.com/songangweb/mcache"
)

type UrlListCache interface {
	Gget(key interface{}) (fullUrl interface{}, ok bool)
	Gadd(key, value interface{}, expirationTime int64)
	Gremove(key interface{})

	Bget(key interface{}) (fullUrl interface{}, ok bool)
	Badd(key, value interface{}, expirationTime int64)
}

type ulServer struct {
	GoodUrlCache *mcache.ARCCache
	BedUrlCache  *mcache.LruCache
}

func NewUrlListCache() UrlListCache {
	return &ulServer{
		GoodUrlCache: GoodUrlCache,
		BedUrlCache:  BedUrlCache,
	}
}

var GoodUrlCache *mcache.ARCCache
var BedUrlCache *mcache.LruCache

type Conf struct {
	GoodUrlLen int
	BedUrlLen  int
}

// InitUrlCache
// 函数名称: InitUrlCache
// 功能: 初始化缓存
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #
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

	GoodUrlCache = goodUrlCache
	BedUrlCache = bedUrlCache
}

func (s *ulServer) Gget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.GoodUrlCache.Get(key)
	return fullUrl, ok
}

func (s *ulServer) Gadd(key, value interface{}, expirationTime int64) {
	s.GoodUrlCache.Add(key, value, expirationTime)
}

func (s *ulServer) Gremove(key interface{}) {
	s.GoodUrlCache.Remove(key)
}

func (s *ulServer) Bget(key interface{}) (fullUrl interface{}, ok bool) {
	fullUrl, _, ok = s.BedUrlCache.Get(key)
	return fullUrl, ok
}

func (s *ulServer) Badd(key, value interface{}, expirationTime int64) {
	s.BedUrlCache.Add(key, value, expirationTime)
}
