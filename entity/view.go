package entity

// SliceView interface - used for providing access to entity ids
type SliceView interface {
	EntityIds() []Id
}

// SetView interface - used for providing access to entity ids
type SetView interface {
	EntityIds() []Id
}
