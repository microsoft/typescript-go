package dirty

import "maps"

type MapBuilder[K comparable, VBase any, VBuilder any] struct {
	base    map[K]VBase
	dirty   map[K]VBuilder
	deleted map[K]struct{}

	toBuilder func(VBase) VBuilder
	build     func(VBuilder) VBase
}

func NewMapBuilder[K comparable, VBase any, VBuilder any](
	base map[K]VBase,
	toBuilder func(VBase) VBuilder,
	build func(VBuilder) VBase,
) *MapBuilder[K, VBase, VBuilder] {
	return &MapBuilder[K, VBase, VBuilder]{
		base:      base,
		dirty:     make(map[K]VBuilder),
		toBuilder: toBuilder,
		build:     build,
	}
}

func (mb *MapBuilder[K, VBase, VBuilder]) Set(key K, value VBuilder) {
	mb.dirty[key] = value
	delete(mb.deleted, key)
}

func (mb *MapBuilder[K, VBase, VBuilder]) Delete(key K) {
	if mb.deleted == nil {
		mb.deleted = make(map[K]struct{})
	}
	mb.deleted[key] = struct{}{}
	delete(mb.dirty, key)
}

func (mb *MapBuilder[K, VBase, VBuilder]) Build() map[K]VBase {
	if len(mb.dirty) == 0 && len(mb.deleted) == 0 {
		return mb.base
	}
	result := maps.Clone(mb.base)
	if result == nil {
		result = make(map[K]VBase)
	}
	for key := range mb.deleted {
		delete(result, key)
	}
	for key, value := range mb.dirty {
		result[key] = mb.build(value)
	}
	return result
}
