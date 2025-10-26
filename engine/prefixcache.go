package engine

import (
	"sync"

	xx "github.com/cespare/xxhash/v2"
)

type KVRef struct{ BlockID int; Tokens int }

type prefixKey struct {
	Model string
	Hash  uint64
}

type PrefixCache struct {
	mu   sync.Mutex
	m    map[prefixKey]KVRef
	tok  Tokenizer
}

func NewPrefixCache() *PrefixCache {
	return &PrefixCache{m: make(map[prefixKey]KVRef), tok: whitespaceTokenizer{}}
}

func NewPrefixCacheWithTokenizer(t Tokenizer) *PrefixCache {
	if t == nil {
		t = whitespaceTokenizer{}
	}
	return &PrefixCache{m: make(map[prefixKey]KVRef), tok: t}
}

// tokenHashes returns running xxhash64 at each token boundary using the configured tokenizer.
func (pc *PrefixCache) tokenHashes(s string) (hashes []uint64, toks []string) {
	toks = pc.tok.Tokenize(s)
	h := xx.New()
	for i, t := range toks {
		if i > 0 {
			h.WriteString(" ")
		}
		h.WriteString(t)
		hashes = append(hashes, h.Sum64())
	}
	return hashes, toks
}

func (pc *PrefixCache) Set(model, prefix string, ref KVRef) {
	hashes, _ := pc.tokenHashes(prefix)
	if len(hashes) == 0 {
		return
	}
	pc.mu.Lock()
	pc.m[prefixKey{Model: model, Hash: hashes[len(hashes)-1]}] = ref
	pc.mu.Unlock()
}

func (pc *PrefixCache) GetLongest(model, full string) (KVRef, string, bool) {
	hashes, toks := pc.tokenHashes(full)
	if len(hashes) == 0 {
		return KVRef{}, "", false
	}
	pc.mu.Lock()
	defer pc.mu.Unlock()
	for i := len(hashes) - 1; i >= 0; i-- {
		if ref, ok := pc.m[prefixKey{Model: model, Hash: hashes[i]}]; ok {
			return ref, pc.tok.Join(toks[:i+1]), true
		}
	}
	return KVRef{}, "", false
}
