package do

type StateMachine[S comparable] struct {
	states []S
	sf     map[S]StateFunc[S]
}

type StateFunc[S comparable] func(pre S) S

func NewStateMachine[S comparable](states []S) *StateMachine[S] {
	return &StateMachine[S]{
		states: Unique(states),
		sf:     make(map[S]StateFunc[S]),
	}
}

func NewStateMachineByMap[S comparable](m map[S]StateFunc[S]) *StateMachine[S] {
	states := make([]S, 0, len(m))
	for k := range m {
		states = append(states, k)
	}

	return &StateMachine[S]{
		states: states,
		sf:     m,
	}
}

func (m *StateMachine[S]) WithFunc(s S, f StateFunc[S]) {
	if !In(m.states, s) {
		return
	}

	if m.sf == nil {
		m.sf = make(map[S]StateFunc[S])
	}

	m.sf[s] = f
}

func (m *StateMachine[S]) Run(s S) S {
	if !In(m.states, s) {
		return s
	}

	return m.sf[s](s)
}
