package hashSimpleLfu

// LFUCache 定义LFUCache接口
type LFUCache interface {

	// Add 向缓存添加一个值。如果已经存在,则更新信息
	Add(key, value *interface{}, expirationTime int64) (ok bool)

	// Get 从缓存中查找一个键的值。
	Get(key *interface{}) (value *interface{}, expirationTime int64, ok bool)

	// Contains 检查某个键是否在缓存中，但不更新缓存的状态
	Contains(key *interface{}) (ok bool)

	// Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态
	Peek(key *interface{}) (value *interface{}, expirationTime int64, ok bool)

	// Remove 从缓存中移除提供的键。
	Remove(key *interface{}) (ok bool)

	// RemoveOldest 从缓存中移除最老的项
	RemoveOldest() (key, value *interface{}, expirationTime int64, ok bool)

	// GetOldest 从缓存中返回最旧的条目
	GetOldest() (key, value *interface{}, expirationTime int64, ok bool)

	// Keys 返回缓存中键的切片，从最老到最新
	Keys() []*interface{}

	// Len 获取缓存已存在的缓存条数
	Len() int

	// Purge 清除所有缓存项
	Purge()

	// PurgeOverdue 清除所有过期缓存项
	PurgeOverdue()

	// Resize 调整缓存大小，返回调整前的数量
	Resize(int) int
}
