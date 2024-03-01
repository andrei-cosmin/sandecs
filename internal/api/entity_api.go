package api

import (
	"github.com/andrei-cosmin/sandata/data"
	"github.com/andrei-cosmin/sandecs/entity"
)

// EntityLinker interface - used for linking entities
type EntityLinker interface {
	Link() entity.Id
	Unlink(entityId entity.Id)
	EntityMask() data.Mask
	GetScheduledRemoves() data.Mask
	Update()
	IsClear() bool
	Refresh()
}
