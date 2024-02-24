package pool

type Pool[T any] struct {
	cursor    int
	container []T
	empty     T
}

func New[T any](size uint) *Pool[T] {
	return &Pool[T]{
		cursor:    -1,
		container: make([]T, size),
	}
}

func (p *Pool[T]) Push(value T) {
	if p.cursor+1 == len(p.container) {
		return
	} else {
		p.cursor++
		p.container[p.cursor] = value
	}
}

func (p *Pool[T]) Pop() (T, bool) {
	if p.cursor == -1 {
		return p.empty, false
	} else {
		value := p.container[p.cursor]
		p.container[p.cursor] = p.empty
		p.cursor--
		return value, true
	}
}
