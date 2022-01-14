# 欢迎使用 mcache 内存缓存包

### mcache是一个基于golang-lru开发的缓存包

mcache 增加了缓存过期时间,增加lfu算法,修改了原有arc算法的依赖结构.
后续还会源源不断增加内存算法.

- 根据过期时间懒汉式删除过期数据,也可主动刷新过期缓存

## why? 为什么要用mcache?

因缓存的使用相关需求,牺牲一部分服务器内存,因减少了网络数据交互,直接使用本机内存,可换取比redis,memcache等更快的缓存速度,
可做为更高一层的缓存需要

## what? 用mcache能做什么?

可作为超高频率数据使用的缓存存储机制

## how? mcache怎么用?

根据需要的不同缓存淘汰算法,使用对应的调用方式

## 代码实现:    
    
    len := 10
    
    // NewLRU 构造一个给定大小的LRU缓存列表
    
    Cache, _ := m_cache.NewLRU(Len)
    
    // NewLFU 构造一个给定大小的LFU缓存列表
    
    Cache, _ := m_cache.NewLFU(Len)
    
    // NewARC 构造一个给定大小的ARC缓存列表
    
    Cache, _ := m_cache.NewARC(Len)
    
    // New2Q 构造一个给定大小的2Q缓存列表
    
    Cache, _ := m_cache.New2Q(Len)
    
    
    // Purge is used to completely clear the cache.
    // Purge 用于完全清除缓存
    
    Cache.Purge()
    
    // PurgeOverdue is used to completely clear the overdue cache.
    // PurgeOverdue 用于清除过期缓存。

    Cache.PurgeOverdue()

    // Add adds a value to the cache. Returns true if an eviction occurred.
    // Add 向缓存添加一个值。如果已经存在,则更新信息
    
    Cache.Add(1,1,1614306658000)
    Cache.Add(2,2,0) // expirationTime 传0代表无过期时间

    // Get looks up a key's value from the cache.
    // Get 从缓存中查找一个键的值

    Cache.Get(2)
    
    // Contains checks if a key is in the cache, without updating the
    // recent-ness or deleting it for being stale.
    // Contains 检查某个键是否在缓存中，但不更新缓存的状态

    Cache.Contains(2)
    
    // Peek returns the key value (or undefined if not found) without updating
    // the "recently used"-ness of the key.
    // Peek 在不更新的情况下返回键值(如果没有找到则返回false),不更新缓存的状态

    Cache.Peek(2)

    // ContainsOrAdd checks if a key is in the cache without updating the
    // recent-ness or deleting it for being stale, and if not, adds the value.
    // Returns whether found and whether an eviction occurred.
    // ContainsOrAdd 判断是否已经存在于缓存中,如果已经存在则不创建及更新内容

    Cache.ContainsOrAdd(3,3,0)
    
    // PeekOrAdd checks if a key is in the cache without updating the
    // recent-ness or deleting it for being stale, and if not, adds the value.
    // Returns whether found and whether an eviction occurred.
    // PeekOrAdd 判断是否已经存在于缓存中,如果已经存在则不更新其键的使用状态

    Cache.PeekOrAdd(4,4,0)

    // Remove removes the provided key from the cache.
    // Remove 从缓存中移除提供的键
    
    Cache.Remove(2)

    // Resize changes the cache size.
    // Resize 调整缓存大小，返回调整前的数量

    len := 100
    Cache.Resize(len)
    
    // RemoveOldest removes the oldest item from the cache.
    // RemoveOldest 从缓存中移除最老的项

    Cache.RemoveOldest()

    // GetOldest returns the oldest entry
    // GetOldest 返回最老的条目
    
    Cache.GetOldest()
    
    // Keys returns a slice of the keys in the cache, from oldest to newest.
    // Keys 返回缓存中键的切片，从最老的到最新的
    
    Cache.Keys()

    // Len returns the number of items in the cache.
    // Len 获取缓存已存在的缓存条数

    Cache.Len()
