package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"reflect"
)

func RegisterComponentLinker[T component.Component](componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	var componentType = reflect.TypeFor[T]().String()
	l := componentLinkManager.(*linkManager)
	return attachLinker(l, componentType, func() api.ComponentLinker {
		return newComponentLinker[T](l.defaultLinkerSize, l.poolCapacity, l.componentIdCursor, componentType, l.entityLinker, l.Set)
	})
}

func RegisterTagLinker(tag component.Tag, componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	l := componentLinkManager.(*linkManager)
	return attachLinker(l, tag, func() api.ComponentLinker {
		return newTagLinker(l.defaultLinkerSize, l.componentIdCursor, tag, l.entityLinker, l.Set)
	})
}

func attachLinker(linkManager *linkManager, label string, constructor func() api.ComponentLinker) api.ComponentLinker {
	if id, ok := linkManager.linkedComponents[label]; ok {
		return linkManager.Get(id)
	} else {
		id = linkManager.componentIdCursor
		instancedLinker := constructor()
		linkManager.linkedComponents[label] = linkManager.componentIdCursor
		linkManager.componentLinkers.Set(
			linkManager.componentIdCursor,
			instancedLinker,
		)
		linkManager.componentIdCursor++
		return instancedLinker
	}
}
