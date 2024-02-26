package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type componentLinker[T component.Component] struct {
	poolCapacity uint
	components   table[T]
	linker
}

func newComponentLinker[T component.Component](size uint, poolCapacity uint, componentId component.Id, componentType string, entityLinker api.EntityContainer, callback func()) api.ComponentLinker {
	var componentTable table[T]
	if poolCapacity > 0 {
		componentTable = newPooledTable[T](size, poolCapacity)
	} else {
		componentTable = newBasicTable[T](size)
	}

	return &componentLinker[T]{
		poolCapacity: poolCapacity,
		components:   componentTable,
		linker:       *newTagLinker(size, componentId, componentType, entityLinker, callback),
	}
}

func (r *componentLinker[T]) Get(entityId entity.Id) *T {
	return r.components.get(entityId)
}

func (r *componentLinker[T]) Link(entityId entity.Id) *T {
	if r.linker.Link(entityId) {
		r.components.set(entityId)
		return r.components.get(entityId)
	}
	return nil
}

func (r *componentLinker[T]) Update(scheduledEntityRemoves *bitset.BitSet) {
	r.linker.Update(scheduledEntityRemoves)
	r.components.clear(r.scheduledRemoves)
}
