package hashSimpleLfu

import (
	"container/heap"
	"time"
)

// entry is used to hold a value in the evictList
// 缓存详细信息
type entry struct {
	key            *interface{} // key
	value          *interface{} // val
	expirationTime int64       // 过期时间
	weight         int         // 访问次数
	index          int         // queue索引
}

// 最小堆
type minHeap []*entry

// Len 堆长度
func (m *minHeap) Len() int {
	return len(*m)
}

// Swap 交换堆元素
func (m *minHeap) Swap(i, j int) {
	// 交换元素
	(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
	// 索引不用交换
	(*m)[i].index = i
	(*m)[j].index = j
}

// Less '<' 是最小堆，'>' 是最大堆  <= 为相同请求次数也会进行交换一次, 保证最近访问的数据,访问次数相同,也会下沉一次
func (m *minHeap) Less(i, j int) bool {
	return (*m)[i].weight <= (*m)[j].weight
}

func (m *minHeap) Push(v interface{}) {
	v.(*entry).index = m.Len()
	*m = append(*m, v.(*entry))
	m.Fix(v.(*entry).index)
}

func (m *minHeap) Pop() interface{} {
	n := m.Len() - 1
	m.Swap(0, n)
	ent := (*m)[n]
	*m = (*m)[:n]
	m.Fix(0)
	return ent
}

func (m *minHeap) Back() interface{} {
	ent := (*m)[0]
	return ent
}

func (m *minHeap) Fix(i int) {
	heap.Fix(m, i)
}

// EvictCallback is used to get a callback when a cache entry is evicted
// EvictCallback 用于在缓存条目被淘汰时的回调函数
type EvictCallback func(key, value *interface{}, expirationTime int64)


// LFU implements a non-thread safe fixed size LFU cache
// LFU 结构体
type LFU struct {

	// 缓存最大条数
	size int

	// 最小堆实现的队列
	minHeap *minHeap

	// 映射缓存
	mapCache map[*interface{}]*entry

	// 回调函数
	onEvict EvictCallback
}

// NewLFU constructs an LFU of the given size
// NewLFU 创建一个新 Cache
func NewLFU(size int, onEvict EvictCallback) (LFUCache, error) {
	minHeap := make(minHeap, 0)
	heap.Init(&minHeap)
	return &LFU{
		size:     size,
		minHeap:  &minHeap,
		mapCache: make(map[*interface{}]*entry),
		onEvict:  onEvict,
	}, nil
}

// Purge is used to completely clear the cache.
// Purge 用于完全清除缓存。
func (l *LFU) Purge() {
	l.mapCache = make(map[*interface{}]*entry)
	minHeap := make(minHeap, 0)
	heap.Init(&minHeap)
	l.minHeap = &minHeap
}

// PurgeOverdue is used to completely clear the overdue cache.
// PurgeOverdue 用于清除过期缓存。
func (l *LFU) PurgeOverdue() {
	for _, v := range l.mapCache {
		// 判断数据是否已经过期
		if checkExpirationTime(v.expirationTime) {
			l.Remove(v.key)
		}
	}
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
// 通过 Add 方法往 Cache 头部增加一个元素，如果存在则更新值
func (l *LFU) Add(key, value *interface{}, expirationTime int64) (ok bool) {

	// 判断此值是否已经超时
	if checkExpirationTime(expirationTime) {
		return false
	}

	// 判断是否已经存在于队列缓存中,则更新数据
	if ent, ok := l.mapCache[key]; ok {
		(*l.minHeap)[ent.index].value = value
		(*l.minHeap)[ent.index].expirationTime = expirationTime
		return ok
	}

	// 如果超出长度，则删除最 '无用' 的元素
	for len(l.mapCache) >= l.size {
		l.RemoveOldest()
	}

	ent := &entry{key: key, value: value, expirationTime: expirationTime}

	l.minHeap.Push(ent)

	l.mapCache[key] = ent

	return true
}

// Get looks up a key's value from the cache.
// Get 从缓存中查找一个键的值。
func (l *LFU) Get(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	// 判读是否存在缓存中,如果存在则,访问次数+1, 调整堆排序
	if ent, ok := l.mapCache[key]; ok {
		// 判断此值是否已经超时,如果超时则进行删除
		if checkExpirationTime((*ent).expirationTime) {
			l.removeEntry(key)
			return nil, 0, false
		}
		ent.weight++
		(*l.minHeap)[ent.index] = ent
		// 进行堆排序
		l.minHeap.Fix(ent.index)
		return ent.value, ent.expirationTime, ok
	}
	return nil, 0, ok
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
// Contains 检查某个键是否在缓存中，但不更新缓存的状态
func (l *LFU) Contains(key *interface{}) (ok bool) {
	// 判读是否存在缓存中
	_, ok = l.mapCache[key]
	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
func (l *LFU) Peek(key *interface{}) (value *interface{}, expirationTime int64, ok bool) {
	ent, ok := l.mapCache[key]
	return ent.value, ent.expirationTime, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
// Remove 从缓存中移除提供的键
func (l *LFU) Remove(key *interface{}) bool {
	if _, ok := l.mapCache[key]; ok {
		l.removeEntry(key)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
// RemoveOldest 从缓存中移除最老的项
func (l *LFU) RemoveOldest() (key, value *interface{}, expirationTime int64, ok bool) {
	if l.Len() == 0 {
		return nil, nil, 0, true
	}
	return l.removeOldest()
}

// GetOldest returns the oldest entry
// GetOldest 返回最老的条目
func (l *LFU) GetOldest() (key, value *interface{}, expirationTime int64, ok bool) {

	// 返回 堆顶元素
	v := l.minHeap.Back()
	if v != nil {

		// 判断此值是否已经超时
		if checkExpirationTime(v.(*entry).expirationTime) {
			return l.GetOldest()
		}

		return v.(*entry).key, v.(*entry).value, v.(*entry).expirationTime, true
	}
	return nil, nil, 0, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Keys 返回缓存的切片，从最老的到最新的。
func (l *LFU) Keys() []*interface{} {
	keys := make([]*interface{}, len(l.mapCache))
	i := 0
	for index, v := range *l.minHeap {

		// 判断此值是否已经超时
		if checkExpirationTime(v.expirationTime) {
			l.Remove(v.key)
		} else {
			keys[index] = v.value
			i++
		}
	}
	return keys
}

// Len returns the number of items in the cache.
// Len 返回缓存中的条数
func (l *LFU) Len() int {
	return len(l.mapCache)
}

// Resize changes the cache size.
// Resize 改变缓存大小。
func (l *LFU) Resize(size int) (evicted int) {
	diff := l.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		l.RemoveOldest()
	}
	l.size = size
	return diff
}

// removeOldest removes the oldest item from the cache.
// removeOldest 从缓存中移除最老的项。
func (l *LFU) removeOldest() (key, value *interface{}, expirationTime int64, ok bool) {
	// 删除堆顶数据
	ent := l.minHeap.Pop()
	// 删除映射缓存
	l.removeMapCache(ent)
	if l.onEvict != nil {
		l.onEvict(ent.(*entry).key, ent.(*entry).value, ent.(*entry).expirationTime)
	}
	return ent.(*entry).key, ent.(*entry).value, ent.(*entry).expirationTime, true
}

// removeMapCache 删除映射缓存
func (l *LFU) removeMapCache(v interface{}) {
	if v == nil {
		return
	}
	delete(l.mapCache, v.(*entry).key)
}

// removeEntry 从缓存中移除一个列表元素
func (l *LFU) removeEntry(key *interface{}) {
	if en, ok := l.mapCache[key]; ok {
		// 数据挪到堆顶
		l.minHeap.Swap(0, en.index)
		// 删除堆顶元素
		v := l.minHeap.Pop()
		// 删除映射缓存
		l.removeMapCache(v)

		if l.onEvict != nil {
			l.onEvict(en.key, en.value, en.expirationTime)
		}
	}
}


// checkExpirationTime is Determine if the cache has expired
// checkExpirationTime 判断缓存是否已经过期
func checkExpirationTime(expirationTime int64) (ok bool) {
	if 0 != expirationTime && expirationTime <= time.Now().UnixNano()/1e6 {
		return true
	}
	return false
}
