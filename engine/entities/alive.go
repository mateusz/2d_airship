package engine

type Alive interface {
	Step(deltaT float64)
}

func (e Entities) Step(deltaT float64) {
	for _, ent := range e {
		a, ok := ent.(Alive)
		if ok {
			a.Step(deltaT)
		}
	}
}
