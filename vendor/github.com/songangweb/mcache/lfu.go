package mcache

import (
	"github.com/songangweb/mcache/simplelfu"
	"sync"
)

// LfuCache is a thread-safe fixed size LRU cache.
// LfuCache 实现一个给定大小的LFU缓存
type LfuCache struct {
	lfu  simplelfu.LFUCache
	lock sync.RWMutex
}

// NewLFU creates an LRU of the given size.
// NewLRU 构造一个给定大小的LRU
func NewLFU(size int) (*LfuCache, error) {
	return NewLfuWithEvict(size, nil)
}

// NewLfuWithEvict constructs a fixed size cache with the given eviction
// callback.
// NewLruWithEvict 用于在缓存条目被淘汰时的回调函数
func NewLfuWithEvict(size int, onEvicted func(key interface{}, value interface{}, expirationTime int64)) (*LfuCache, error) {
	lfu, _ := simplelfu.NewLFU(size, simplelfu.EvictCallback(onEvicted))
	c := &LfuCache{
		lfu: lfu,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
// Purge 用于完全清除缓存
func (c *LfuCache) Purge() {
	c.lock.Lock()
	c.lfu.Purge()
	c.lock.Unlock()
}

// PurgeOverdue is used to completely clear the overdue cache.
// PurgeOverdue 用于清除过期缓存。
func (c *LfuCache) PurgeOverdue() {
	c.lock.Lock()
	c.lfu.PurgeOverdue()
	c.lock.Unlock()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
// Add 向缓存添加一个值。如果已经存在,则更新信息
func (c *LfuCache) Add(key, value interface{}, expirationTime int64) (evicted bool) {
	c.lock.Lock()
	evicted = c.lfu.Add(key, value, expirationTime)
	c.lock.Unlock()
	return evicted
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值
func (c *LfuCache) Get(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	value, expirationTime, ok = c.lfu.Get(key)
	c.lock.Unlock()
	return value, expirationTime, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (c *LfuCache) Contains(key interface{}) bool {
	c.lock.RLock()
	containKey := c.lfu.Contains(key)
	c.lock.RUnlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (c *LfuCache) Peek(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.RLock()
	value, expirationTime, ok = c.lfu.Peek(key)
	c.lock.RUnlock()
	return value, expirationTime, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// ContainsOrAdd 判断是否已经存在于缓存中,如果已经存在则不创建及更新内容
func (c *LfuCache) ContainsOrAdd(key, value interface{}, expirationTime int64) (ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lfu.Contains(key) {
		return true, false
	}
	evicted = c.lfu.Add(key, value, expirationTime)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// PeekOrAdd 判断是否已经存在于缓存中,如果已经存在则不更新其键的使用状态
func (c *LfuCache) PeekOrAdd(key, value interface{}, expirationTime int64) (previous interface{}, ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	previous, expirationTime, ok = c.lfu.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lfu.Add(key, value, expirationTime)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
// Remove 从缓存中移除提供的键
func (c *LfuCache) Remove(key interface{}) (present bool) {
	c.lock.Lock()
	present = c.lfu.Remove(key)
	c.lock.Unlock()
	return
}

// Resize changes the cache size.
// Resize 调整缓存大小，返回调整前的数量
func (c *LfuCache) Resize(size int) (evicted int) {
	c.lock.Lock()
	evicted = c.lfu.Resize(size)
	c.lock.Unlock()
	return evicted
}

// RemoveOldest removes the oldest item from the cache.
// RemoveOldest 从缓存中移除最老的项
func (c *LfuCache) RemoveOldest() (key interface{}, value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	key, value, expirationTime, ok = c.lfu.RemoveOldest()
	c.lock.Unlock()
	return
}

// GetOldest returns the oldest entry
// GetOldest 返回最老的条目
func (c *LfuCache) GetOldest() (key interface{}, value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	key, value, expirationTime, ok = c.lfu.GetOldest()
	c.lock.Unlock()
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Keys 返回缓存中键的切片，从最老的到最新的
func (c *LfuCache) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lfu.Keys()
	c.lock.RUnlock()
	return keys
}

// Len returns the number of items in the cache.
// Len 获取缓存已存在的缓存条数
func (c *LfuCache) Len() int {
	c.lock.RLock()
	length := c.lfu.Len()
	c.lock.RUnlock()
	return length
}
