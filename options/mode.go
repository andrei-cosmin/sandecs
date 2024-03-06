package options

type Mode byte

const (
	Standard Mode = iota
	Pooled
	Compact
)

const (
	DefaultNumEntities   = 128
	DefaultNumComponents = 16
	DefaultPoolCapacity  = 1024
)
