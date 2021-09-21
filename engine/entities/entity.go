package engine

type Entity interface {
}

type Entities []Entity

func NewEntities() Entities {
	return make([]Entity, 0)
}

func (e Entities) Add(ent Entity) Entities {
	return append(e, ent)
}

func (e Entities) Remove(ent Entity) {
	for i, ml := range e {
		if ent == ml {
			copy(e[i:], e[i+1:])
			e = e[:len(e)-1]
			break
		}
	}
}

func (e Entities) Len() int {
	return len(e)
}
