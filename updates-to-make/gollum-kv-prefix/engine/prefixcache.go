package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type KVRef struct{ BlockID int; Tokens int }

type PrefixCache struct{
	mu sync.Mutex
	m  map[string]KVRef
}
func NewPrefixCache()*PrefixCache{ return &PrefixCache{m: make(map[string]KVRef)} }
func (pc *PrefixCache) key(model, prefix string) string {
	h := sha256.Sum256([]byte(model + "|" + prefix))
	return hex.EncodeToString(h[:])
}
func (pc *PrefixCache) Get(model, prefix string)(KVRef,bool){
	k := pc.key(model,prefix); pc.mu.Lock(); defer pc.mu.Unlock(); v,ok := pc.m[k]; return v,ok
}
func (pc *PrefixCache) Set(model, prefix string, ref KVRef){
	k := pc.key(model,prefix); pc.mu.Lock(); defer pc.mu.Unlock(); pc.m[k] = ref
}
