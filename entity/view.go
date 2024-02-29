package entity

// View interface - used for providing access to entity ids
type View interface {
	EntityIds() []Id
}
