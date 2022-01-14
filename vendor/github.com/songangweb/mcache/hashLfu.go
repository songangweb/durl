package mcache

import (
	"crypto/md5"
	hashSimpleLfu "github.com/songangweb/mcache/hashsimplelfu"
	"math"
	"runtime"
	"sync"
)

// HashLfuCache is a thread-safe fixed size LRU cache.
// HashLfuCache 实现一个给定大小的LFU缓存
type HashLfuCache struct {
	list     []*HashLfuCacheOne
	sliceNum int
	size     int
}

type HashLfuCacheOne struct {
	lfu  hashSimpleLfu.LFUCache
	lock sync.RWMutex
}


// NewHashLFU creates an LRU of the given size.
// NewHashLFU 构造一个给定大小的LRU
func NewHashLFU(size, sliceNum int) (*HashLfuCache, error) {
	return NewHashLfuWithEvict(size, sliceNum, nil)
}

// NewHashLfuWithEvict constructs a fixed size cache with the given eviction
// callback.
// NewHashLfuWithEvict 用于在缓存条目被淘汰时的回调函数
func NewHashLfuWithEvict(size, sliceNum int, onEvicted func(key, value *interface{}, expirationTime int64)) (*HashLfuCache, error) {
	if 0 == sliceNum {
		// 设置为当前cpu数量
		sliceNum = runtime.NumCPU()
	}
	if size < sliceNum {
		size = sliceNum
	}

	// 计算出每个分片的数据长度
	lfuLen := int(math.Ceil(float64(size/sliceNum)))

	var h HashLfuCache
	h.size = size
	h.sliceNum = sliceNum
	h.list = make([]*HashLfuCacheOne, sliceNum)
	for i := 0; i < sliceNum; i++ {
		l, _ := hashSimpleLfu.NewLFU(lfuLen, onEvicted)
		h.list[i] = &HashLfuCacheOne{
			lfu: l,
		}
	}

	return &h, nil
}


// Add adds a value to the cache. Returns true if an eviction occurred.
// Add 向缓存添加一个值。如果已经存在,则更新信息
func (h *HashLfuCache) Add(key, value *interface{}, expirationTime int64) (evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	evicted = h.list[sliceKey].lfu.Add(key, value, expirationTime)
	h.list[sliceKey].lock.Unlock()
	return evicted
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值
func (h *HashLfuCache) Get(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	value, expirationTime, ok = h.list[sliceKey].lfu.Get(key)
	h.list[sliceKey].lock.Unlock()
	return value, expirationTime, ok
}

// Purge is used to completely clear the cache.
// Purge 用于完全清除缓存
func (h *HashLfuCache) Purge() {
	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		h.list[i].lfu.Purge()
		h.list[i].lock.Unlock()
	}
}

// PurgeOverdue is used to completely clear the overdue cache.
// PurgeOverdue 用于清除过期缓存。
func (h *HashLfuCache) PurgeOverdue() {
	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		h.list[i].lfu.PurgeOverdue()
		h.list[i].lock.Unlock()
	}
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (h *HashLfuCache) Contains(key *interface{}) bool {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.RLock()
	containKey := h.list[sliceKey].lfu.Contains(key)
	h.list[sliceKey].lock.RUnlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (h *HashLfuCache) Peek(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.RLock()
	value, expirationTime, ok = h.list[sliceKey].lfu.Peek(key)
	h.list[sliceKey].lock.RUnlock()
	return value, expirationTime, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// ContainsOrAdd 判断是否已经存在于缓存中,如果已经存在则不创建及更新内容
func (h *HashLfuCache) ContainsOrAdd(key, value *interface{}, expirationTime int64) (ok, evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	defer h.list[sliceKey].lock.Unlock()

	if h.list[sliceKey].lfu.Contains(key) {
		return true, false
	}
	evicted = h.list[sliceKey].lfu.Add(key, value, expirationTime)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
// PeekOrAdd 判断是否已经存在于缓存中,如果已经存在则不更新其键的使用状态
func (h *HashLfuCache) PeekOrAdd(key, value *interface{}, expirationTime int64) (previous interface{}, ok, evicted bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	defer h.list[sliceKey].lock.Unlock()

	previous, expirationTime, ok = h.list[sliceKey].lfu.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = h.list[sliceKey].lfu.Add(key, value, expirationTime)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
// Remove 从缓存中移除提供的键
func (h *HashLfuCache) Remove(key *interface{}) (present bool) {
	sliceKey := h.qumo(key)

	h.list[sliceKey].lock.Lock()
	present = h.list[sliceKey].lfu.Remove(key)
	h.list[sliceKey].lock.Unlock()
	return
}

// Resize changes the cache size.
// Resize 调整缓存大小，返回调整前的数量
func (h *HashLfuCache) Resize(size int) (evicted int) {
	if size < h.sliceNum {
		size = h.sliceNum
	}

	// 计算出每个分片的数据长度
	lfuLen := int(math.Ceil(float64(size/h.sliceNum)))

	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.Lock()
		evicted = h.list[i].lfu.Resize(lfuLen)
		h.list[i].lock.Unlock()
	}

	return evicted
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Keys 返回缓存中键的切片，从最老的到最新的
func (h *HashLfuCache) Keys() []*interface{} {
	sliceKeys := make([][]*interface{}, h.sliceNum)

	sliceMaxSize := 0
	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.RLock()
		sliceKeys[i] = h.list[i].lfu.Keys()
		// 获取各个切片中最大的条数
		sliceSize := h.list[i].lfu.Len()
		if sliceSize > sliceMaxSize {
			sliceMaxSize = sliceSize
		}
		h.list[i].lock.RUnlock()
	}

	keys := make([]*interface{}, h.Len())
	k := 0
	for i := 0; i < sliceMaxSize; i++ {
		for j := 0; j < h.sliceNum; j++ {
			if i < len(sliceKeys[j]) {
				keys[k] = sliceKeys[j][i]
				k++
			}
		}
	}
	return keys
}


// Len returns the number of items in the cache.
// Len 获取缓存已存在的缓存条数
func (h *HashLfuCache) Len() int {
	var length = 0

	for i := 0; i < h.sliceNum; i++ {
		h.list[i].lock.RLock()
		length = length + h.list[i].lfu.Len()
		h.list[i].lock.RUnlock()
	}
	return length
}

func (h *HashLfuCache)qumo(key *interface{}) int {
	str := InterfaceToString(*key)
	return int(md5.Sum([]byte(str))[0]) % h.sliceNum
}
