package engine

import (
	"strings"
	"sync"

	xx "github.com/cespare/xxhash/v2"
)

type KVRef struct{ BlockID int; Tokens int }

// We store entries keyed by (model, hash64). The hash corresponds to the
// whitespace-tokenized prefix up to some token boundary. Collisions are
// acceptable for this cache tier.
type prefixKey struct {
	Model string
	Hash  uint64
}

type PrefixCache struct {
	mu sync.Mutex
	m  map[prefixKey]KVRef
}

func NewPrefixCache() *PrefixCache { return &PrefixCache{m: make(map[prefixKey]KVRef)} }

// tokenHashes returns the running xxhash64 at each token boundary.
// e.g. toks = ["foo","bar","baz"] -> hashes for: "foo", "foo bar", "foo bar baz"
func tokenHashes(s string) (hashes []uint64, toks []string) {
	toks = strings.Fields(s)
	h := xx.New()
	first := true
	for _, t := range toks {
		if !first {
			h.WriteString(" ")
		}
		h.WriteString(t)
		first = false
		hashes = append(hashes, h.Sum64())
	}
	return hashes, toks
}

// Set registers a KVRef for the FULL prefix string (model+prefix).
// We store only the final hash for that prefix.
func (pc *PrefixCache) Set(model, prefix string, ref KVRef) {
	hashes, toks := tokenHashes(prefix)
	if len(hashes) == 0 {
		return
	}
	pc.mu.Lock()
	pc.m[prefixKey{Model: model, Hash: hashes[len(hashes)-1]}] = ref
	pc.mu.Unlock()

	_ = toks // (kept to clarify semantics; may be useful for future partial storage)
}

// GetLongest returns the KVRef and the matched prefix text for the LONGEST
// cached prefix of `full`. O(n) membership checks, descending in length.
func (pc *PrefixCache) GetLongest(model, full string) (KVRef, string, bool) {
	hashes, toks := tokenHashes(full)
	if len(hashes) == 0 {
		return KVRef{}, "", false
	}
	pc.mu.Lock()
	defer pc.mu.Unlock()
	for i := len(hashes) - 1; i >= 0; i-- {
		if ref, ok := pc.m[prefixKey{Model: model, Hash: hashes[i]}]; ok {
			return ref, strings.Join(toks[:i+1], " "), true
		}
	}
	return KVRef{}, "", false
}
