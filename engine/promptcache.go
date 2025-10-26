package engine

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type CacheEntry struct {
	Key    string
	Tokens []string
}

type PromptCache struct {
	mu    sync.Mutex
	cap   int
	ll    *list.List
	items map[string]*list.Element
}

type node struct{ e CacheEntry }

func NewPromptCache(capacity int) *PromptCache {
	if capacity <= 0 {
		capacity = 1024
	}
	return &PromptCache{cap: capacity, ll: list.New(), items: make(map[string]*list.Element)}
}

func (pc *PromptCache) key(prompt string, model string, temp float32, maxTokens int) string {
	h := sha256.Sum256([]byte(prompt + "|" + model + "|" + string(rune(int(temp*100))) + "|" + string(rune(maxTokens))))
	return hex.EncodeToString(h[:])
}

func (pc *PromptCache) Get(prompt, model string, temp float32, maxTokens int) (tokens []string, ok bool) {
	k := pc.key(prompt, model, temp, maxTokens)
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if el, ok := pc.items[k]; ok {
		pc.ll.MoveToFront(el)
		n := el.Value.(*node)
		return n.e.Tokens, true
	}
	return nil, false
}

func (pc *PromptCache) Put(prompt, model string, temp float32, maxTokens int, tokens []string) {
	k := pc.key(prompt, model, temp, maxTokens)
	pc.mu.Lock()
	defer pc.mu.Unlock()
	if el, ok := pc.items[k]; ok {
		n := el.Value.(*node)
		n.e.Tokens = tokens
		pc.ll.MoveToFront(el)
		return
	}
	el := pc.ll.PushFront(&node{e: CacheEntry{Key: k, Tokens: tokens}})
	pc.items[k] = el
	if pc.ll.Len() > pc.cap {
		old := pc.ll.Back()
		if old != nil {
			pc.ll.Remove(old)
			delete(pc.items, old.Value.(*node).e.Key)
		}
	}
}
