package entity

import "github.com/andrei-cosmin/sandata/bit"

// View provides access to filtered entity IDs and their bitmask.
type View interface {
	EntityIds() []Id
	EntityMask() bit.Mask
}

// MaskView provides access to an entity bitmask.
type MaskView interface {
	EntityMask() bit.Mask
}
