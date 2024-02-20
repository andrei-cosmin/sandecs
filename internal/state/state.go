package state

type State struct {
	isUpdated bool
}

func New() State {
	return State{isUpdated: true}
}

func (l *State) Mark() {
	l.isUpdated = false
}

func (l *State) Reset() {
	l.isUpdated = true
}

func (l *State) IsUpdated() bool {
	return l.isUpdated
}
