package component

type Component interface{}

type Linker[T Component] interface {
	Link(entity uint, component T)
	Get(entity uint) *T
	Has(entity uint) bool
	Remove(entity uint)
}
