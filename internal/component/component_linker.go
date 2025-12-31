package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/entity"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"github.com/andrei-cosmin/sandecs/options"
	"github.com/bits-and-blooms/bitset"
)

// componentLinker manages component instances of type T.
type componentLinker[T component.Component] struct {
	baseLinker
	poolCapacity uint
	components   table[T]
	additions    *bitset.BitSet
	onLink       func(*T)
	onUnlink     func(*T)
}

func newComponentLinker[T component.Component](
	mode options.Mode,
	size, poolCapacity uint,
	componentId component.Id, componentType string,
	entityLinker entity.MaskView,
	callback func(),
) api.ComponentLinker {
	if poolCapacity <= 0 {
		poolCapacity = options.DefaultPoolCapacity
	}
	if size <= 0 {
		size = options.DefaultNumComponents
	}

	var componentTable table[T]
	switch mode {
	default:
	case options.Standard:
		componentTable = newBasicTable[T](size)
	case options.Pooled:
		componentTable = newPooledTable[T](size, poolCapacity)
	case options.Compact:
		componentTable = newCompactTable[T](size)
	}

	return &componentLinker[T]{
		poolCapacity: poolCapacity,
		components:   componentTable,
		additions:    bitset.New(size),
		baseLinker:   *newBaseLinker(size, componentId, componentType, entityLinker, callback),
	}
}

// Get returns the component for the entity, or nil if not linked.
func (r *componentLinker[T]) Get(entityId entity.Id) *T {
	return r.components.get(entityId)
}

// Link attaches a component to the entity and returns it. Returns nil if already linked.
func (r *componentLinker[T]) Link(entityId entity.Id) *T {
	if r.baseLinker.Link(entityId) {
		r.components.set(entityId)
		r.additions.Set(entityId)
		return r.components.get(entityId)
	}
	return nil
}

// SetLinkHook sets a callback invoked when a component is linked.
func (r *componentLinker[T]) SetLinkHook(onLink func(*T)) {
	r.onLink = onLink
}

// SetUnlinkHook sets a callback invoked when a component is unlinked.
func (r *componentLinker[T]) SetUnlinkHook(onUnlink func(*T)) {
	r.onUnlink = onUnlink
}

// RemoveLinkHook clears the link hook.
func (r *componentLinker[T]) RemoveLinkHook(onLink func()) {
	r.onLink = nil
}

// RemoveUnlinkHook clears the unlink hook.
func (r *componentLinker[T]) RemoveUnlinkHook(onUnlink func()) {
	r.onUnlink = nil
}

// CleanScheduledInstances processes scheduled removals and triggers hooks.
func (r *componentLinker[T]) CleanScheduledInstances() {
	if r.onLink != nil {
		for addedEntityId, hasNext := r.additions.NextSet(0); hasNext; addedEntityId, hasNext = r.additions.NextSet(addedEntityId + 1) {
			r.onLink(r.components.get(addedEntityId))
		}
	}
	r.components.clear(r.scheduledRemoves, r.onUnlink)
	r.additions.ClearAll()
}
