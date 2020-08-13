package storage

// Local implements local storage as an in-memory hetrogenous list of entites.
type Local struct {
	Entities []interface{}
}

// Query local storage.
func (l *Local) Query(plist []Predicate) (interface{}, bool) {
	for _, ent := range l.Entities {
		if ok := Predicates(plist).Apply(ent); ok {
			return ent, ok
		}
	}
	return nil, false
}

// Save to local storage.
func (l *Local) Save(ent interface{}) error {
	l.Entities = append(l.Entities, ent)
	return nil
}
