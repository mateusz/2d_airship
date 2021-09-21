package engine

import (
	"sort"
)

type Positionable interface {
	GetZ() float64
	GetX() float64
	GetY() float64
}

func (e Entities) ByZ() Entities {
	eWork := NewEntities()
	eWork = append(eWork, e...)
	sort.SliceStable(eWork, func(i, j int) bool {
		pi, oki := eWork[i].(Positionable)
		pj, okj := eWork[j].(Positionable)

		orderi := 0.0
		orderj := 0.0
		if oki {
			orderi = pi.GetZ()
		}
		if okj {
			orderj = pj.GetZ()
		}

		return orderi < orderj
	})
	return eWork
}

func (e Entities) ByReverseZ() Entities {
	eWork := NewEntities()
	eWork = append(eWork, e...)
	sort.SliceStable(eWork, func(i, j int) bool {
		pi, oki := eWork[i].(Positionable)
		pj, okj := eWork[j].(Positionable)

		orderi := 0.0
		orderj := 0.0
		if oki {
			orderi = pi.GetZ()
		}
		if okj {
			orderj = pj.GetZ()
		}

		return orderi > orderj
	})
	return eWork
}
