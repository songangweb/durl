package tool

import (
	"strconv"
	"strings"
)

type Trie struct {
	children [256]*Trie
}

func Constructor() Trie {
	return Trie{}
}

func (t *Trie) Add(ip string) {
	byteList := ipv4ToByte(ip)
	cur := t
	for _, c := range byteList {
		if cur.children[c] == nil {
			cur.children[c] = &Trie{}
		}
		cur = cur.children[c]
	}
}

func (t *Trie) Search(ip string) bool {
	byteList := ipv4ToByte(ip)
	cur := t
	for _, c := range byteList {
		if cur.children[c] == nil {
			return false
		}
		cur = cur.children[c]
	}
	return true
}

func ipv4ToByte(ipAddr string) []byte {
	if ipAddr == "::1" {
		ipAddr = "127.0.0.1"
	}
	bits := strings.Split(ipAddr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var b []byte
	b = append(b, byte(uint32(b0)))
	b = append(b, byte(uint32(b1)))
	b = append(b, byte(uint32(b2)))
	b = append(b, byte(uint32(b3)))

	return b
}
