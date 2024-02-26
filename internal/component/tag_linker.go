package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/bits-and-blooms/bitset"
)

type linker struct {
	componentId      component.Id
	componentType    string
	entityLinker     api.EntityContainer
	callback         func()
	scheduledRemoves *bitset.BitSet
	linkedEntities   *bitset.BitSet
}

func newTagLinker(size uint, componentId component.Id, componentType string, entityLinker api.EntityContainer, callback func()) *linker {
	return &linker{
		componentId:      componentId,
		componentType:    componentType,
		entityLinker:     entityLinker,
		callback:         callback,
		scheduledRemoves: bitset.New(size),
		linkedEntities:   bitset.New(size),
	}
}

func (r *linker) Link(entityId entity.Id) bool {
	if !r.entityLinker.EntityIds().Test(entityId) || r.linkedEntities.Test(entityId) {
		return false
	}

	r.linkedEntities.Set(entityId)
	r.callback()
	return true
}

func (r *linker) Has(entityId entity.Id) bool {
	return r.linkedEntities.Test(entityId)
}

func (r *linker) Unlink(entityId entity.Id) {
	if !r.Has(entityId) {
		return
	}
	if r.scheduledRemoves.Test(entityId) {
		return
	}

	r.scheduledRemoves.Set(entityId)
	r.callback()
}

func (r *linker) ComponentId() component.Id {
	return r.componentId
}

func (r *linker) EntityIds() *bitset.BitSet {
	return r.linkedEntities
}

func (r *linker) Update(scheduledEntityRemoves *bitset.BitSet) {
	r.scheduledRemoves.InPlaceUnion(scheduledEntityRemoves)
	r.linkedEntities.InPlaceDifference(r.scheduledRemoves)
}

func (r *linker) Refresh() {
	r.scheduledRemoves.ClearAll()
}
