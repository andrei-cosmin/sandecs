package component

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/api"
	"github.com/andrei-cosmin/hakkt/internal/sparse"
	"github.com/bits-and-blooms/bitset"
)

type linker[T component.Component] struct {
	componentId      component.Id
	componentType    string
	callback         func()
	components       *sparse.Array[*T]
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

func newLinker[T component.Component](size uint, componentId component.Id, componentType string, callback func()) api.ComponentLinker {
	return &linker[T]{
		componentId:      componentId,
		componentType:    componentType,
		callback:         callback,
		components:       sparse.New[*T](size),
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

func (r *linker[T]) Link(entityId entity.Id) {
	if r.linkedEntities.Test(entityId) {
		return
	}

	r.linkedEntities.Set(entityId)
	r.components.Set(entityId, new(T))
	r.callback()
}

func (r *linker[T]) Get(entityId entity.Id) *T {
	return r.components.Get(entityId)
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
	r.components.Clear(r.scheduledRemoves)
	r.linkedEntities.InPlaceDifference(r.scheduledRemoves)
	r.scheduledRemoves.ClearAll()
}
