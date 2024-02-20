package component

import (
	"github.com/andrei-cosmin/hakkt/component"
	"github.com/andrei-cosmin/hakkt/internal/sparse"
	"github.com/bits-and-blooms/bitset"
)

type Linker struct {
	componentId      uint
	componentType    string
	callback         func()
	components       *sparse.Array[*component.Component]
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

func newLinker(size uint, componentId uint, componentType string, callback func()) *Linker {
	return &Linker{
		componentId:      componentId,
		componentType:    componentType,
		callback:         callback,
		components:       sparse.New[*component.Component](size),
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

func (r *Linker) Link(entity uint, component *component.Component) {
	r.linkedEntities.Set(entity)
	r.components.Set(entity, component)
	r.callback()
}

func (r *Linker) Get(entity uint) *component.Component {
	return r.components.Get(entity)
}

func (r *Linker) Has(entity uint) bool {
	return r.linkedEntities.Test(entity)
}

func (r *Linker) Remove(entity uint) {
	if !r.Has(entity) {
		return
	}

	if r.scheduledRemoves.Test(entity) {
		return
	}
	r.scheduledRemoves.Set(entity)
	r.callback()
}

func (r *Linker) GetComponentId() uint {
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
