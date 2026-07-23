package ast

import (
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/core"
)

type (
	NodeId   uint64
	SymbolId uint64
)

// Atomic ids

var (
	nextNodeId   atomic.Uint64
	nextSymbolId atomic.Uint64
)

func GetNodeId(node *Node) NodeId {
	return NodeId(getId(&node.id, &nextNodeId))
}

func GetSymbolId(symbol *Symbol) SymbolId {
	return SymbolId(getId(&symbol.id, &nextSymbolId))
}

func getId(id *atomic.Uint64, counter *atomic.Uint64) uint64 {
	value := id.Load()
	if value == 0 {
		// Worst case, we burn a few ids if we have to CAS.
		value = counter.Add(1)
		if !id.CompareAndSwap(0, value) {
			value = id.Load()
		}
	}
	return value
}

// IDs are handed out in page-aligned blocks so that a link store page only ever holds IDs minted
// by one allocator. Blocks grow from one page to maxIdBlockSize as an allocator keeps allocating.
const (
	minIdBlockSize = core.LinkStorePageSize
	maxIdBlockSize = 16 * core.LinkStorePageSize
)

// IdAllocator assigns node and symbol IDs from contiguous blocks reserved from the global counters,
// which keeps the IDs one client mints clustered in the same link store pages. An IdAllocator must
// not be used concurrently.
type IdAllocator struct {
	nodeIds   idBlock
	symbolIds idBlock
}

func (a *IdAllocator) GetNodeId(node *Node) NodeId {
	if id := node.id.Load(); id != 0 {
		return NodeId(id)
	}
	return NodeId(a.nodeIds.assign(&node.id, &nextNodeId))
}

func (a *IdAllocator) GetSymbolId(symbol *Symbol) SymbolId {
	if id := symbol.id.Load(); id != 0 {
		return SymbolId(id)
	}
	return SymbolId(a.symbolIds.assign(&symbol.id, &nextSymbolId))
}

type idBlock struct {
	next uint64
	end  uint64
	size uint64
}

func (b *idBlock) assign(id *atomic.Uint64, counter *atomic.Uint64) uint64 {
	if b.next == b.end {
		b.reserve(counter)
	}
	value := b.next
	if !id.CompareAndSwap(0, value) {
		return id.Load()
	}
	b.next++
	return value
}

// reserve claims the next page-aligned block of IDs from counter.
func (b *idBlock) reserve(counter *atomic.Uint64) {
	b.size = min(max(2*b.size, minIdBlockSize), maxIdBlockSize)
	for {
		last := counter.Load()
		first := (last/core.LinkStorePageSize + 1) * core.LinkStorePageSize
		if counter.CompareAndSwap(last, first+b.size-1) {
			b.next = first
			b.end = first + b.size
			return
		}
	}
}
