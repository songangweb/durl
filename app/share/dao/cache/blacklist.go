package cache

import (
	"durl/app/share/tool"
)

var Blacklist *blacklist

type blacklist struct {
	trie *tool.Trie
}

func InitBlacklist() {
	ipTrie := tool.Constructor()
	Blacklist = &blacklist{
		trie: &ipTrie,
	}
}

func (b *blacklist) Add(ip string) {
	b.trie.Add(ip)
}

func (b *blacklist) Search(ip string) bool {
	return b.trie.Search(ip)
}
