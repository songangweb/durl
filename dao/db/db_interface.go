package db

// 定义LFUCache接口
type Db interface {

	// Add 向缓存添加一个值。如果已经存在,则更新信息
	Add(key, value interface{}, expirationTime int64) (ok bool)
}
