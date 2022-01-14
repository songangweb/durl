package mcache

import (
	"github.com/songangweb/mcache/simplelfu"
	"github.com/songangweb/mcache/simplelru"
	"sync"
)

// ARCCache is a thread-safe fixed size Adaptive Replacement LfuCache (ARC).
// ARC is an enhancement over the standard LRU cache in that tracks both
// frequency and recency of use. This avoids a burst in access to new
// entries from evicting the frequently used older entries. It adds some
// additional tracking overhead to a standard LRU cache, computationally
// it is roughly 2x the cost, and the extra memory overhead is linear
// with the size of the cache. ARC has been patented by IBM, but is
// similar to the TwoQueueCache (2Q) which requires setting parameters.
// ARCCache 是一个线程安全的固定大小自适应替换缓存(ARC)。
// ARC是对标准LRU缓存的一个增强，它可以同时跟踪这两个缓存
// 使用的频率和频率。这避免了访问新内容的突然爆发
// 删除常用的旧条目。它增加了一些
// 额外的跟踪开销到一个标准的LRU缓存，计算
// 大约是开销的2倍，额外的内存开销是线性的
// 使用缓存的大小。ARC已经被IBM申请了专利，但它是
// 类似于TwoQueueCache (2Q)，需要设置参数。
type ARCCache struct {
	// Size为缓存的总容量
	size int
	// P是对T1或T2的动态偏好
	p int

	t1 simplelru.LRUCache // T1 is the LRU for recently accessed items
	b1 simplelru.LRUCache // B1 is the LRU for evictions from t1

	t2 simplelfu.LFUCache // T2 is the LFU for frequently accessed items
	b2 simplelfu.LFUCache // B2 is the LFU for evictions from t2

	lock sync.RWMutex
}

// NewARC creates an ARC of the given size
func NewARC(size int) (*ARCCache, error) {
	// Create the sub LRUs
	t1, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	b1, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	t2, err := simplelfu.NewLFU(size, nil)
	if err != nil {
		return nil, err
	}
	b2, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	// Initialize the ARC
	c := &ARCCache{
		size: size,
		p:    0,
		t1:   t1,
		b1:   b1,
		t2:   t2,
		b2:   b2,
	}
	return c, nil
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值。
func (c *ARCCache) Get(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the value is contained in T1 (recent), then
	// promote it to T2 (frequent)
	if val, expirationTime, ok := c.t1.Peek(key); ok {
		c.t1.Remove(key)
		c.t2.Add(key, val, expirationTime)
		return val, expirationTime, ok
	}

	// Check if the value is contained in T2 (frequent)
	if val, expirationTime, ok := c.t2.Get(key); ok {
		return val, expirationTime, ok
	}

	// No hit
	return nil, expirationTime, false
}

// Add adds a value to the cache.
// Add 向缓存添加一个值。如果已经存在,则更新信息
func (c *ARCCache) Add(key, value interface{}, expirationTime int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the value is contained in T1 (recent), and potentially
	// promote it to frequent T2
	if c.t1.Contains(key) {
		c.t1.Remove(key)
		c.t2.Add(key, value, expirationTime)
		return
	}

	// Check if the value is already in T2 (frequent) and update it
	if c.t2.Contains(key) {
		c.t2.Add(key, value, expirationTime)
		return
	}

	// Check if this value was recently evicted as part of the
	// recently used list
	if c.b1.Contains(key) {
		// T1 set is too small, increase P appropriately
		delta := 1
		b1Len := c.b1.Len()
		b2Len := c.b2.Len()
		if b2Len > b1Len {
			delta = b2Len / b1Len
		}
		if c.p+delta >= c.size {
			c.p = c.size
		} else {
			c.p += delta
		}

		// Potentially need to make room in the cache
		if c.t1.Len()+c.t2.Len() >= c.size {
			c.replace(false)
		}

		// Remove from B1
		c.b1.Remove(key)

		// Add the key to the frequently used list
		c.t2.Add(key, value, expirationTime)
		return
	}

	// Check if this value was recently evicted as part of the
	// frequently used list
	if c.b2.Contains(key) {
		// T2 set is too small, decrease P appropriately
		delta := 1
		b1Len := c.b1.Len()
		b2Len := c.b2.Len()
		if b1Len > b2Len {
			delta = b1Len / b2Len
		}
		if delta >= c.p {
			c.p = 0
		} else {
			c.p -= delta
		}

		// Potentially need to make room in the cache
		if c.t1.Len()+c.t2.Len() >= c.size {
			c.replace(true)
		}

		// Remove from B2
		c.b2.Remove(key)

		// Add the key to the frequently used list
		c.t2.Add(key, value, expirationTime)
		return
	}

	// Potentially need to make room in the cache
	if c.t1.Len()+c.t2.Len() >= c.size {
		c.replace(false)
	}

	// Keep the size of the ghost buffers trim
	if c.b1.Len() > c.size-c.p {
		c.b1.RemoveOldest()
	}
	if c.b2.Len() > c.p {
		c.b2.RemoveOldest()
	}

	// Add to the recently seen list
	c.t1.Add(key, value, expirationTime)
	return
}

// replace is used to adaptively evict from either T1 or T2
// based on the current learned value of P
// replace 用于自适应地从T1或T2中驱逐,根据P的当前学习值
func (c *ARCCache) replace(b2ContainsKey bool) {
	t1Len := c.t1.Len()
	if t1Len > 0 && (t1Len > c.p || (t1Len == c.p && b2ContainsKey)) {
		k, _, expirationTime, ok := c.t1.RemoveOldest()
		if ok {
			c.b1.Add(k, nil, expirationTime)
		}
	} else {
		k, _, expirationTime, ok := c.t2.RemoveOldest()
		if ok {
			c.b2.Add(k, nil, expirationTime)
		}
	}
}

// Len returns the number of cached entries
// Len 获取缓存已存在的缓存条数
func (c *ARCCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.t1.Len() + c.t2.Len()
}

// Keys returns all the cached keys
// Keys 返回缓存中键的切片，从最老到最新
func (c *ARCCache) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	k1 := c.t1.Keys()
	k2 := c.t2.Keys()
	return append(k1, k2...)
}

// Remove is used to purge a key from the cache
// Remove 从缓存中移除提供的键。
func (c *ARCCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.t1.Remove(key) {
		return
	}
	if c.t2.Remove(key) {
		return
	}
	if c.b1.Remove(key) {
		return
	}
	if c.b2.Remove(key) {
		return
	}
}

// Purge is used to clear the cache
// Purge 清除所有缓存项
func (c *ARCCache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.t1.Purge()
	c.t2.Purge()
	c.b1.Purge()
	c.b2.Purge()
}

// Contains is used to check if the cache contains a key
// without updating recency or frequency.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (c *ARCCache) Contains(key interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.t1.Contains(key) || c.t2.Contains(key)
}

// Peek is used to inspect the cache value of a key
// without updating recency or frequency.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (c *ARCCache) Peek(key interface{}) (value interface{}, expirationTime int64, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if val, expirationTime, ok := c.t1.Peek(key); ok {
		return val, expirationTime, ok
	}
	return c.t2.Peek(key)
}
