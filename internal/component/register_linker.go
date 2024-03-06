package component

import (
	"github.com/andrei-cosmin/sandecs/component"
	"github.com/andrei-cosmin/sandecs/internal/api"
	"reflect"
)

// RegisterComponentLinker method - registers a component linker with the component link manager
func RegisterComponentLinker[T component.Component](componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	var componentType = reflect.TypeFor[T]().String()
	l := componentLinkManager.(*linkManager)
	// Register the component linker for the given parametrized type T
	return registerLinker(l, componentType, func() api.ComponentLinker {
		return newComponentLinker[T](l.mode, l.defaultLinkerSize, l.poolCapacity, l.componentIdCursor, componentType, l.entityLinker, l.Set)
	})
}

// RegisterTagLinker method -  registers a tagLinker with the component link manager
func RegisterTagLinker(tag component.Tag, componentLinkManager api.ComponentLinkManager) api.ComponentLinker {
	l := componentLinkManager.(*linkManager)
	// Register the tagLinker for the given tag
	return registerLinker(l, tag, func() api.ComponentLinker {
		return newTagLinker(l.defaultLinkerSize, l.componentIdCursor, tag, l.entityLinker, l.Set)
	})
}

// registerLinker method - attaches a new tagLinker to the link manager for a given type, or returns the existing tagLinker
func registerLinker(linkManager *linkManager, label string, constructor func() api.ComponentLinker) api.ComponentLinker {
	if id, ok := linkManager.linkedComponents[label]; ok {
		// If the tagLinker already exists, return it
		return linkManager.Get(id)
	} else {
		// Get the next available component id
		id = linkManager.componentIdCursor

		// Create a new tagLinker and attach it to the link manager, using the given type (label / component type)
		instancedLinker := constructor()
		linkManager.linkedComponents[label] = linkManager.componentIdCursor

		// Set the tagLinker in the component linkers array
		linkManager.componentLinkers.Set(
			linkManager.componentIdCursor,
			instancedLinker,
		)

		// Increment the component id cursor, for the next available component id
		linkManager.componentIdCursor++

		// Return the instanced tagLinker
		return instancedLinker
	}
}
