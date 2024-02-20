package component

import (
	"github.com/andrei-cosmin/hakkt-ecs/internal/sparse"
	"github.com/andrei-cosmin/hakkt-ecs/internal/state"
	"github.com/bits-and-blooms/bitset"
)

type LinkManager struct {
	linkedComponents  map[string]uint
	Linkers           *sparse.Array[*Linker]
	componentIdCursor uint
	state.State
}

func NewLinkManager(size uint) *LinkManager {
	return &LinkManager{
		linkedComponents:  make(map[string]uint),
		Linkers:           sparse.New[*Linker](size),
		componentIdCursor: 0,
		State:             state.New(),
	}
}

func (l *LinkManager) Link(componentType string) *Linker {
	var componentId uint

	if id, ok := l.linkedComponents[componentType]; ok {
		componentId = id
	} else {
		componentId = l.componentIdCursor
		l.linkedComponents[componentType] = l.componentIdCursor
		l.Linkers.Set(
			l.componentIdCursor,
			newLinker(l.Linkers.Size(), l.componentIdCursor, componentType, l.Mark),
		)
		l.componentIdCursor++
	}

	return l.Linkers.Get(componentId)
}

func (l *LinkManager) Find(componentId uint) *Linker {
	return l.Linkers.Get(componentId)
}

func (l *LinkManager) UpdateLinks(scheduledEntityRemoves *bitset.BitSet) {
	l.Reset()
	for index := uint(0); index < l.componentIdCursor; index++ {
		var resolver = l.Linkers.Get(index)
		if resolver != nil {
			resolver.update(scheduledEntityRemoves)
		}
	}
}
