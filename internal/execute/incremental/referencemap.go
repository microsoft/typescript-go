package incremental

import (
	"maps"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type referenceMap struct {
	references       map[tspath.Path][]tspath.Path
	referencedBy     map[tspath.Path][]tspath.Path
	referencedByOnce sync.Once
}

func (r *referenceMap) makeReferences(size int) {
	if size != 0 {
		r.references = make(map[tspath.Path][]tspath.Path, size)
	}
}

func (r *referenceMap) storeReferences(path tspath.Path, refs []tspath.Path) {
	r.references[path] = refs
}

func (r *referenceMap) getReferences(path tspath.Path) []tspath.Path {
	return r.references[path]
}

func (r *referenceMap) getPathsWithReferences() []tspath.Path {
	return slices.Collect(maps.Keys(r.references))
}

func (r *referenceMap) getReferencedBy(path tspath.Path) []tspath.Path {
	r.referencedByOnce.Do(func() {
		referencedBy := make(map[tspath.Path][]tspath.Path)
		for key, value := range r.references {
			for _, ref := range value {
				referencedBy[ref] = append(referencedBy[ref], key)
			}
		}
		r.referencedBy = make(map[tspath.Path][]tspath.Path, len(referencedBy))
		maps.Copy(r.referencedBy, referencedBy)
	})
	return r.referencedBy[path]
}
