package benchmarks

const (
	numPosition         = 90000
	numPositionVelocity = 10000
	numEdits            = 10000
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
