package engine

import "container/list"

type KVBlock struct {
	ID       int
	Tokens   int
	Pinned   bool
	KVHandle int // Bound metal KV handle for this block (0 = unbound). One handle per block.
}

type KVPager struct {
	blocks []*KVBlock
	free   *list.List
	lru    *list.List
	nextID int
}

func NewKVPager() *KVPager {
	p := &KVPager{free: list.New(), lru: list.New()}
	for i := 0; i < 64; i++ {
		b := &KVBlock{ID: i, Tokens: 4096}
		p.blocks = append(p.blocks, b)
		p.free.PushBack(b)
	}
	p.nextID = len(p.blocks)
	return p
}

func (p *KVPager) Allocate(tokens int) *KVBlock {
	for el := p.free.Front(); el != nil; el = el.Next() {
		b := el.Value.(*KVBlock)
		if b.Tokens >= tokens {
			p.free.Remove(el)
			p.lru.PushFront(b)
			return b
		}
	}
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

func (p *KVPager) Free(b *KVBlock) {
	if b == nil {
		return
	}
	for el := p.lru.Front(); el != nil; el = el.Next() {
		if el.Value.(*KVBlock).ID == b.ID {
			p.lru.Remove(el)
			break
		}
	}
	b.Pinned = false
	p.free.PushBack(b)
}

func (p *KVPager) Pin(b *KVBlock) {
	if b != nil {
		b.Pinned = true
	}
}

func (p *KVPager) Unpin(b *KVBlock) {
	if b != nil {
		b.Pinned = false
	}
}

func (p *KVPager) ByID(id int) *KVBlock {
	for _, b := range p.blocks {
		if b.ID == id {
			return b
		}
	}
	return nil
}
