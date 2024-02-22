package component

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/entity"
	"github.com/andrei-cosmin/hakkt/internal/sparse"
	"github.com/bits-and-blooms/bitset"
)

type Linker struct {
	componentId      component.Id
	componentType    string
	callback         func()
	components       *sparse.Array[*component.Component]
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

func newLinker(size uint, componentId component.Id, componentType string, callback func()) *Linker {
	return &Linker{
		componentId:      componentId,
		componentType:    componentType,
		callback:         callback,
		components:       sparse.New[*component.Component](size),
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

func (r *Linker) Link(entityId entity.Id, component *component.Component) {
	r.linkedEntities.Set(entityId)
	r.components.Set(entityId, component)
	r.callback()
}

func (r *Linker) Get(entityId entity.Id) *component.Component {
	return r.components.Get(entityId)
}

func (r *Linker) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

func (r *Linker) Remove(entityId entity.Id) {
	if !r.Has(entityId) {
		return
	}

	if r.scheduledRemoves.Test(entityId) {
		return
	}
	r.scheduledRemoves.Set(entityId)
	r.callback()
}

func (r *Linker) GetComponentId() component.Id {
	return r.componentId
}

func (r *Linker) GetEntities() *bitset.BitSet {
	return r.linkedEntities
}

func (r *Linker) update(scheduledEntityRemoves *bitset.BitSet) {
	r.scheduledRemoves.InPlaceUnion(scheduledEntityRemoves)
	r.components.Clear(r.scheduledRemoves)
	r.linkedEntities.InPlaceDifference(r.scheduledRemoves)
	r.scheduledRemoves.ClearAll()
}
