package benchmarks

const (
	numPosition         = 9000
	numPositionVelocity = 1000
	numEdits            = 1000
)

type position struct {
	X float64
	Y float64
}
type velocity struct {
	X float64
	Y float64
}

type name struct {
	Name string
}
