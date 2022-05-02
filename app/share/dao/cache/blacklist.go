package cache

import (
	"sync"

	"durl/app/share/tool"
)

var (
	Blacklist         *tool.Trie
	BlacklistConnLock sync.RWMutex
)

// InitBlacklist
// 函数名称: InitBlacklist
// 功能: 初始化黑名单 ip树接口
// 输入参数:
// 输出参数:
// 返回: 返回请求结果
// 实现描述:
// 注意事项:
// 作者: # ang.song # 2021-11-17 15:15:42 #
func InitBlacklist() *tool.Trie {
	ipTrie := tool.Constructor()
	return &ipTrie
}
