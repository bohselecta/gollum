package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
)

type KVRef struct {
	BlockID int
	Tokens  int
}

type PrefixCache struct {
	mu sync.Mutex
	// key: model + "|" + hex(sha256(prefix))
	m map[string]KVRef
}

func NewPrefixCache() *PrefixCache {
	return &PrefixCache{m: make(map[string]KVRef)}
}

func key(model, prefix string) string {
	h := sha256.Sum256([]byte(prefix))
	return model + "|" + hex.EncodeToString(h[:])
}

// Set stores an exact prefix mapping.
func (pc *PrefixCache) Set(model, prefix string, ref KVRef) {
	pc.mu.Lock()
	pc.m[key(model, prefix)] = ref
	pc.mu.Unlock()
}

// GetExact is still available if you need it.
func (pc *PrefixCache) GetExact(model, prefix string) (KVRef, bool) {
	pc.mu.Lock()
	ref, ok := pc.m[key(model, prefix)]
	pc.mu.Unlock()
	return ref, ok
}

// GetLongest returns the longest cached prefix for `full` and its text that matched.
// If none, ok=false.
func (pc *PrefixCache) GetLongest(model, full string) (ref KVRef, matchedPrefix string, ok bool) {
	// Tokenize naïvely by whitespace to avoid mid-word splits.
	// Swap with your real tokenizer later.
	toks := strings.Fields(full)
	if len(toks) == 0 {
		return KVRef{}, "", false
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	// Try longest → shortest
	for n := len(toks); n >= 1; n-- {
		prefix := strings.Join(toks[:n], " ")
		if r, ex := pc.m[key(model, prefix)]; ex {
			return r, prefix, true
		}
	}
	return KVRef{}, "", false
}
