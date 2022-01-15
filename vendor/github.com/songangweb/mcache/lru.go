package mcache

import (
	"github.com/songangweb/mcache/simplelru"
	"sync"
)

// LruCache is a thread-safe fixed size LRU cache.
// LruCache 实现一个给定大小的LRU缓存
type LruCache struct {
	lru  simplelru.LRUCache
	lock sync.RWMutex
}

// NewLRU creates an LRU of the given size.
// NewLRU 构造一个给定大小的LRU
func NewLRU(size int) (*LruCache, error) {
	return NewLruWithEvict(size, nil)
}

// NewLruWithEvict constructs a fixed size cache with the given eviction
// callback.
// NewLruWithEvict 用于在缓存条目被淘汰时的回调函数
func NewLruWithEvict(size int, onEvicted func(key interface{}, value interface{}, expirationTime int64)) (*LruCache, error) {
	lru, err := simplelru.NewLRU(size, simplelru.EvictCallback(onEvicted))
	if err != nil {
		return nil, err
	}
	c := &LruCache{
		lru: lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
// Purge 清除所有缓存项
func (c *LruCache) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// PurgeOverdue is used to completely clear the overdue cache.
// PurgeOverdue 用于清除过期缓存。
func (c *LruCache) PurgeOverdue() {
	c.lock.Lock()
	c.lru.PurgeOverdue()
	c.lock.Unlock()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
// Add 向缓存添加一个值。如果已经存在,则更新信息
func (c *LruCache) Add(key, value interface{}, expirationTime int64) (evicted bool) {
	c.lock.Lock()
	evicted = c.lru.Add(key, value, expirationTime)
	c.lock.Unlock()
	return evicted
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值。
func (c *LruCache) Get(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	value, expirationTime, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, expirationTime, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (c *LruCache) Contains(key interface{}) bool {
	c.lock.RLock()
	containKey := c.lru.Contains(key)
	c.lock.RUnlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (c *LruCache) Peek(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.RLock()
	value, expirationTime, ok = c.lru.Peek(key)
	c.lock.RUnlock()
	return value, expirationTime, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// ContainsOrAdd 检查键是否在缓存中，而不更新
// 最近或删除它，因为它是陈旧的，如果不是，添加值。
// 返回是否找到和是否发生了驱逐。
func (c *LruCache) ContainsOrAdd(key, value interface{}, expirationTime int64) (ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value, expirationTime)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// PeekOrAdd 如果一个key在缓存中，那么这个key就不会被更新
// 最近或删除它，因为它是陈旧的，如果不是，添加值。
// 返回是否找到和是否发生了驱逐。
func (c *LruCache) PeekOrAdd(key, value interface{}, expirationTime int64) (previous interface{}, ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	previous, expirationTime, ok = c.lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lru.Add(key, value, expirationTime)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
// Remove 从缓存中移除提供的键。
func (c *LruCache) Remove(key interface{}) (present bool) {
	c.lock.Lock()
	present = c.lru.Remove(key)
	c.lock.Unlock()
	return
}

// Resize changes the cache size.
// Resize 调整缓存大小，返回调整前的数量
func (c *LruCache) Resize(size int) (evicted int) {
	c.lock.Lock()
	evicted = c.lru.Resize(size)
	c.lock.Unlock()
	return evicted
}

// RemoveOldest removes the oldest item from the cache.
// RemoveOldest 从缓存中移除最老的项
func (c *LruCache) RemoveOldest() (key interface{}, value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	key, value, expirationTime, ok = c.lru.RemoveOldest()
	c.lock.Unlock()
	return
}

// GetOldest returns the oldest entry
// GetOldest 从缓存中返回最旧的条目
func (c *LruCache) GetOldest() (key interface{}, value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	key, value, expirationTime, ok = c.lru.GetOldest()
	c.lock.Unlock()
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Keys 返回缓存中键的切片，从最老到最新
func (c *LruCache) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lru.Keys()
	c.lock.RUnlock()
	return keys
}

// Len returns the number of items in the cache.
// Len 获取缓存已存在的缓存条数
func (c *LruCache) Len() int {
	c.lock.RLock()
	length := c.lru.Len()
	c.lock.RUnlock()
	return length
}
