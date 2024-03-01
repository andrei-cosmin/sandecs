package entity

import "github.com/andrei-cosmin/sandata/data"

// View interface - used for providing access to entity ids
type View interface {
	EntityIds() []Id
	EntityMask() data.Mask
}

type MaskView interface {
	EntityMask() data.Mask
}
