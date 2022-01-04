package cache

import (
	"durl/app/share/tool"
)

type BlackListCache interface {
	Add(ip string)
	Search(ip string) bool
}

func NewBlackListCache() BlackListCache {
	return &blcServer{
		Blacklist: Blacklist,
	}
}

type blcServer struct {
	Blacklist *tool.Trie
}

var Blacklist *tool.Trie

func InitBlacklist() {
	ipTrie := tool.Constructor()
	Blacklist = &ipTrie
}

func (b *blcServer) Add(ip string) {
	b.Blacklist.Add(ip)
}

func (b *blcServer) Search(ip string) bool {
	return b.Blacklist.Search(ip)
}
