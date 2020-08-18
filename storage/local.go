package storage

// Local implements local storage as an in-memory hetrogenous list of entites.
type Local struct {
	Entities []interface{}
}

// Query the first entity that satisfies the filters.
func (l *Local) Query(filters ...func(interface{}) bool) (interface{}, bool) {
	for _, ent := range l.Entities {
		if Apply(ent, filters...) {
			return ent, true
		}
	}
	return nil, false
}

// List all entities that satisfy the filters.
func (l *Local) List(filters ...func(interface{}) bool) []interface{} {
	var list []interface{}
	for _, ent := range l.Entities {
		if Apply(ent, filters...) {
			list = append(list, ent)
		}
	}
	return list
}

// Save to local storage.
func (l *Local) Save(ent interface{}) error {
	l.Entities = append(l.Entities, ent)
	return nil
}
