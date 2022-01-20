package cache

import (
	"sync"

	"durl/app/share/tool"
)

var (
	Blacklist         *tool.Trie
	BlacklistConnLock sync.RWMutex
)

func InitBlacklist() *tool.Trie {
	ipTrie := tool.Constructor()
	return &ipTrie
}
