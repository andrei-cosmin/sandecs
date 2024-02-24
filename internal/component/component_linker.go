package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type linker[T component.Component] struct {
	poolCapacity     uint
	componentId      component.Id
	componentType    string
	entityLinker     api.EntityContainer
	callback         func()
	components       table[T]
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

func newLinker[T component.Component](size uint, poolCapacity uint, componentId component.Id, componentType string, entityLinker api.EntityContainer, callback func()) api.ComponentLinker {
	var table table[T]
	if poolCapacity > 0 {
		table = newPooledTable[T](size, poolCapacity)
	} else {
		table = newBasicTable[T](size)
	}

	return &linker[T]{
		poolCapacity:     poolCapacity,
		componentId:      componentId,
		componentType:    componentType,
		entityLinker:     entityLinker,
		callback:         callback,
		components:       table,
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

func (r *linker[T]) Link(entityId entity.Id) *T {
	if !r.entityLinker.EntityIds().Test(entityId) || r.linkedEntities.Test(entityId) {
		return nil
	}

	r.linkedEntities.Set(entityId)
	r.components.set(entityId)
	r.callback()

	return r.components.get(entityId)
}

func (r *linker[T]) Get(entityId entity.Id) *T {
	return r.components.get(entityId)
}

func (r *linker[T]) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

func (r *linker[T]) Unlink(entityId entity.Id) {
	if !r.Has(entityId) {
		return
	}
	if r.scheduledRemoves.Test(entityId) {
		return
	}

	r.scheduledRemoves.Set(entityId)
	r.callback()
}

func (r *linker[T]) ComponentId() component.Id {
	return r.componentId
}

func (r *linker[T]) EntityIds() *bitset.BitSet {
	return r.linkedEntities
}

func (r *linker[T]) Update(scheduledEntityRemoves *bitset.BitSet) {
	r.scheduledRemoves.InPlaceUnion(scheduledEntityRemoves)
	r.components.clear(r.scheduledRemoves)
	r.linkedEntities.InPlaceDifference(r.scheduledRemoves)
	r.scheduledRemoves.ClearAll()
}
