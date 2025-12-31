package options

// Mode defines the component storage strategy.
type Mode byte

// Storage modes.
const (
	Standard Mode = iota // Flat array storage
	Pooled               // Instance reuse via pooling
	Compact              // Dense storage via sparse set
)

// Default sandbox configuration values.
const (
	DefaultNumEntities   = 128
	DefaultNumComponents = 16
	DefaultPoolCapacity  = 1024
)
