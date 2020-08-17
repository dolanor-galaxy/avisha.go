package storage

// Storer provides persistence for entities.
type Storer interface {
	Query(predicates []Predicate) (interface{}, bool)
	List(predicates []Predicate) []interface{}
	Save(ent Entity) error
	Delete(ent Entity) error
}

// Entity is a unique object that has an ID.
// Used to compare for equality when saving objects.
type Entity interface {
	ID() string
}

// Predicate selects an an entity based on some criteria.
// Predicate is responsible for type asserting the entity.
type Predicate interface {
	Apply(ent interface{}) bool
}

// PredicateFunc is a standalone func that implements Predicate.
type PredicateFunc func(ent interface{}) bool

// Apply the predicate func.
func (fn PredicateFunc) Apply(ent interface{}) bool {
	return fn(ent)
}

// Predicates treats a list of predicates as one big predicate.
type Predicates []Predicate

// Apply each predicate to the entity.
// If any of the predicates fail, this predicate as a whole fails.
func (predicates Predicates) Apply(ent interface{}) bool {
	for _, p := range predicates {
		if ok := p.Apply(ent); !ok {
			return false
		}
	}
	return true
}
