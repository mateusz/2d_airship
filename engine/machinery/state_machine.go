package machinery

type Machine struct {
	States  []*State
	Current *State
}

type State struct {
	Name    string
	Ts      []Transition
	OnEntry func()
}

type Transition struct {
	Check func() bool
	To    *State
}

func (m *Machine) Run() {
	m.Current.Run(m)
}

func (s *State) Run(m *Machine) {
	for _, t := range s.Ts {
		if t.Check() {
			m.Current = t.To
			m.Current.OnEntry()
		}
	}
}

func NewTwoStateLatch(entry1, entry2 func(), t12, t21 func() bool) Machine {
	s1 := &State{
		Name: "One",
		Ts: []Transition{
			{Check: t12},
		},
	}
	s2 := &State{
		Name: "Two",
		Ts: []Transition{
			{Check: t21, To: s1},
		},
	}
	s1.Ts[0].To = s2
	m := Machine{
		States:  []*State{s1, s2},
		Current: s1,
	}

	return m
}
