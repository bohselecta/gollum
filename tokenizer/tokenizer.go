package tokenizer

import (
	"strings"
	"sync"
)

// Tokenizer interface for pluggable tokenization
type Tokenizer interface {
	Encode(text string) ([]int, error)
	Decode(ids []int) string
	VocabSize() int
}

// SimpleBytePairEncoding is a minimal BPE implementation
// For production, use a library like github.com/klauspost/go-bpe
type SimpleBytePairEncoding struct {
	vocab     map[string]int
	vocabRev  map[int]string
	mu        sync.RWMutex
	vocabSize int
}

// NewSimpleBPE creates a simple BPE tokenizer
func NewSimpleBPE(vocabSize int) *SimpleBytePairEncoding {
	return &SimpleBytePairEncoding{
		vocab:     make(map[string]int),
		vocabRev:  make(map[int]string),
		vocabSize: vocabSize,
	}
}

func (bpe *SimpleBytePairEncoding) Encode(text string) ([]int, error) {
	// For now, implement simple tokenization
	// In production, load vocab from GGUF metadata and implement BPE merge rules
	
	bpe.mu.RLock()
	defer bpe.mu.RUnlock()
	
	// Simple fallback: split on whitespace and assign sequential IDs
	words := strings.Fields(text)
	ids := make([]int, 0, len(words))
	
	for _, word := range words {
		// Assign a hash-based ID (consistent mapping)
		id := hashString(word) % int32(bpe.vocabSize)
		
		// Add bos/eos tokens
		ids = append(ids, int(id))
	}
	
	// Add EOS token (simplified)
	ids = append(ids, 2)
	
	return ids, nil
}

func (bpe *SimpleBytePairEncoding) Decode(ids []int) string {
	bpe.mu.RLock()
	defer bpe.mu.RUnlock()
	
	// Simple reverse mapping
	var tokens []string
	for _, id := range ids {
		if id == 2 { // EOS
			break
		}
		// In production, look up in vocabRev
		tokens = append(tokens, string(rune(id)))
	}
	
	return strings.Join(tokens, " ")
}

func (bpe *SimpleBytePairEncoding) VocabSize() int {
	return bpe.vocabSize
}

// hashString creates a consistent hash for a string
func hashString(s string) int32 {
	var hash int32
	for _, c := range s {
		hash = hash*31 + c
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}

// WhitespaceTokenizer is the fallback used in prefixcache.go
type WhitespaceTokenizer struct{}

func (w WhitespaceTokenizer) Tokenize(s string) []string {
	return strings.Fields(s)
}

func (w WhitespaceTokenizer) Join(tokens []string) string {
	return strings.Join(tokens, " ")
}
