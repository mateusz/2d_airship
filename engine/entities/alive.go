package engine

type Alive interface {
	Step(deltaT float64)
}
