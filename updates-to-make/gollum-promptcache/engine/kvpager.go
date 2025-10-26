package engine

import "container/list"

// KVBlock represents a contiguous slice of KV memory (abstract handle).
type KVBlock struct {
	ID      int
	Tokens  int   // capacity in tokens
	Pinned  bool
	LastUse int64 // logical clock
}

// KVPager manages blocks with a simple LRU (not wired to actual kernels yet).
type KVPager struct {
	blocks []*KVBlock
	free   *list.List
	lru    *list.List
	nextID int
}

func NewKVPager() *KVPager {
	p := &KVPager{free: list.New(), lru: list.New()}
	// Pre-create some blocks to simulate availability.
	for i:=0;i<64;i++ {
		b := &KVBlock{ID: i, Tokens: 2048}
		p.blocks = append(p.blocks, b)
		p.free.PushBack(b)
	}
	p.nextID = len(p.blocks)
	return p
}

func (p *KVPager) Allocate(tokens int) *KVBlock {
	// naive first-fit
	for el := p.free.Front(); el != nil; el = el.Next() {
		b := el.Value.(*KVBlock)
		if b.Tokens >= tokens {
			p.free.Remove(el)
			p.lru.PushFront(b)
			return b
		}
	}
	// no block: evict LRU non-pinned
	for el := p.lru.Back(); el != nil; el = el.Prev() {
		b := el.Value.(*KVBlock)
		if !b.Pinned {
			p.lru.Remove(el)
			p.lru.PushFront(b)
			return b
		}
	}
	return nil
}

func (p *KVPager) Touch(b *KVBlock) {
	if b == nil { return }
	for el := p.lru.Front(); el != nil; el = el.Next() {
		if el.Value.(*KVBlock).ID == b.ID {
			p.lru.MoveToFront(el); return
		}
	}
}

func (p *KVPager) Free(b *KVBlock) {
	if b == nil { return }
	for el := p.lru.Front(); el != nil; el = el.Next() {
		if el.Value.(*KVBlock).ID == b.ID {
			p.lru.Remove(el); break
		}
	}
	p.free.PushBack(b)
}
