package component

import (
	"reflect"

	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
)

// RegisterComponentLinker registers a component linker for type T.
func RegisterComponentLinker[T component.Component](componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	componentType := reflect.TypeFor[T]().String()
	l := componentLinkManager.(*linkManager)
	return registerLinker(l, componentType, func() api.ComponentLinker {
		return newComponentLinker[T](l.mode, l.defaultLinkerSize, l.poolCapacity, l.componentIdCursor, componentType, l.entityLinker, l.Set)
	})
}

// RegisterTagLinker registers a tag linker.
func RegisterTagLinker(tag component.Tag, componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	l := componentLinkManager.(*linkManager)
	return registerLinker(l, tag, func() api.ComponentLinker {
		return newTagLinker(l.defaultLinkerSize, l.componentIdCursor, tag, l.entityLinker, l.Set)
	})
}

func registerLinker(linkManager *linkManager, label string, constructor func() api.ComponentLinker) api.ComponentLinker {
	if id, ok := linkManager.linkedComponents[label]; ok {
		return linkManager.Get(id)
	}
	id := linkManager.componentIdCursor
	instancedLinker := constructor()
	linkManager.linkedComponents[label] = id
	linkManager.componentLinkers.Set(id, instancedLinker)
	linkManager.componentIdCursor++
	return instancedLinker
}
