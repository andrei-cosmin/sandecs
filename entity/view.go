package entity

import "github.com/andrei-cosmin/sandata/bit"

// View interface - used for providing access to entity ids
type View interface {
	EntityIds() []Id
	EntityMask() bit.Mask
}

// MaskView interface - used for providing access to entity masks
type MaskView interface {
	EntityMask() bit.Mask
}
