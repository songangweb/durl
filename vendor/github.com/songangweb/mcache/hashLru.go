package mcache

import (
	"crypto/md5"
	hashSimpleLru "github.com/songangweb/mcache/hashsimplelru"
	"math"
	"runtime"
	"sync"
)

// HashLruCache is a thread-safe fixed size LRU cache.
// HashLruCache 实现一个给定大小的LRU缓存
type HashLruCache struct {
	list     []*HashLruCacheOne
	sliceNum int
	size     int
}

type HashLruCacheOne struct {
	lru  hashSimpleLru.LRUCache
	lock sync.RWMutex
}

// NewHashLRU creates an LRU of the given size.
// NewHashLRU 构造一个给定大小的LRU
func NewHashLRU(size, sliceNum int) (*HashLruCache, error) {
	return NewHashLruWithEvict(size, sliceNum, nil)
}

// NewHashLruWithEvict constructs a fixed size cache with the given eviction
// callback.
// NewHashLruWithEvict 用于在缓存条目被淘汰时的回调函数
func NewHashLruWithEvict(size, sliceNum int, onEvicted func(key, value *interface{}, expirationTime int64)) (*HashLruCache, error) {
	if 0 == sliceNum {
		// 设置为当前cpu数量
		sliceNum = runtime.NumCPU()
	}
	if size < sliceNum {
		size = sliceNum
	}

	// 计算出每个分片的数据长度
	lruLen := int(math.Ceil(float64(size/sliceNum)))
	var h HashLruCache
	h.size = size
	h.sliceNum = sliceNum
	h.list = make([]*HashLruCacheOne, sliceNum)
	for i := 0; i < sliceNum; i++ {
		l, _ := hashSimpleLru.NewLRU(lruLen, onEvicted)
		h.list[i] = &HashLruCacheOne{
			lru: l,
		}
	}

	return &h, nil
}

// Purge is used to completely clear the cache.
// Purge 清除所有缓存项
func (h *HashLruCache) Purge() {
	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		h.list[i].lru.Purge()
		h.list[i].lock.Unlock()
	}
}

// PurgeOverdue is used to completely clear the overdue cache.
// PurgeOverdue 用于清除过期缓存。
func (h *HashLruCache) PurgeOverdue() {
	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		h.list[i].lru.PurgeOverdue()
		h.list[i].lock.Unlock()
	}
}

// Add adds a value to the cache. Returns true if an eviction occurred.
// Add 向缓存添加一个值。如果已经存在,则更新信息
func (h *HashLruCache) Add(key, value *interface{}, expirationTime int64) (evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	evicted = h.list[sliceKey].lru.Add(key, value, expirationTime)
	h.list[sliceKey].lock.Unlock()
	return evicted
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值。
func (h *HashLruCache) Get(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	value, expirationTime, ok = h.list[sliceKey].lru.Get(key)
	h.list[sliceKey].lock.Unlock()
	return value, expirationTime, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (h *HashLruCache) Contains(key *interface{}) bool {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.RLock()
	containKey := h.list[sliceKey].lru.Contains(key)
	h.list[sliceKey].lock.RUnlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (h *HashLruCache) Peek(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.RLock()
	value, expirationTime, ok = h.list[sliceKey].lru.Peek(key)
	h.list[sliceKey].lock.RUnlock()
	return value, expirationTime, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// ContainsOrAdd 检查键是否在缓存中，而不更新
// 最近或删除它，因为它是陈旧的，如果不是，添加值。
// 返回是否找到和是否发生了驱逐。
func (h *HashLruCache) ContainsOrAdd(key, value *interface{}, expirationTime int64) (ok, evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	defer h.list[sliceKey].lock.Unlock()

	if h.list[sliceKey].lru.Contains(key) {
		return true, false
	}
	evicted = h.list[sliceKey].lru.Add(key, value, expirationTime)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// PeekOrAdd 如果一个key在缓存中，那么这个key就不会被更新
// 最近或删除它，因为它是陈旧的，如果不是，添加值。
// 返回是否找到和是否发生了驱逐。
func (h *HashLruCache) PeekOrAdd(key, value *interface{}, expirationTime int64) (previous interface{}, ok, evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	defer h.list[sliceKey].lock.Unlock()

	previous, expirationTime, ok = h.list[sliceKey].lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = h.list[sliceKey].lru.Add(key, value, expirationTime)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
// Remove 从缓存中移除提供的键。
func (h *HashLruCache) Remove(key *interface{}) (present bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	present = h.list[sliceKey].lru.Remove(key)
	h.list[sliceKey].lock.Unlock()
	return
}

// Resize changes the cache size.
// Resize 调整缓存大小，返回调整前的数量
func (h *HashLruCache) Resize(size int) (evicted int) {
	if size < h.sliceNum {
		size = h.sliceNum
	}

	// 计算出每个分片的数据长度
	lruLen := int(math.Ceil(float64(size/h.sliceNum)))

	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		evicted = h.list[i].lru.Resize(lruLen)
		h.list[i].lock.Unlock()
	}
	return evicted
}

//// Keys returns a slice of the keys in the cache, from oldest to newest.
//// Keys 返回缓存中键的切片，从最老到最新
//func (h *HashLruCache) Keys() []*interface{} {
//	h.lock.RLock()
//	keys := h.lru.Keys()
//	h.lock.RUnlock()
//	return keys
//}

// Len returns the number of items in the cache.
// Len 获取缓存已存在的缓存条数
func (h *HashLruCache) Len() int {
	var length = 0

	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.RLock()
		length = length + h.list[i].lru.Len()
		h.list[i].lock.RUnlock()
	}
	return length
}

func (h *HashLruCache)qumo(key *interface{}) int {
	str := InterfaceToString(*key)
	return int(md5.Sum([]byte(str))[0]) % h.sliceNum
}
