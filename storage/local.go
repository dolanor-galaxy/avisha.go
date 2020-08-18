package storage

import "fmt"

// Local implements local storage as an in-memory hetrogenous list of entites.
type Local struct {
	Entities []Entity
}

// Query the first entity that satisfies the filters.
func (l *Local) Query(filters ...func(Entity) bool) (Entity, bool) {
	for _, ent := range l.Entities {
		if Apply(ent, filters...) {
			return ent, true
		}
	}
	return nil, false
}

// List all entities that satisfy the filters.
func (l *Local) List(filters ...func(Entity) bool) []Entity {
	var list []Entity
	for _, ent := range l.Entities {
		if Apply(ent, filters...) {
			list = append(list, ent)
		}
	}
	return list
}

// Create saves a new entity to local storage.
func (l *Local) Create(ent Entity) error {
	for _, entity := range l.Entities {
		if entity.ID() == ent.ID() {
			return fmt.Errorf("already exists")
		}
	}
	l.Entities = append(l.Entities, ent)
	return nil
}

// Update saves an existing entity to local storage.
func (l *Local) Update(ent Entity) error {
	for ii, entity := range l.Entities {
		if ent.ID() == entity.ID() {
			l.Entities[ii] = ent
			return nil
		}
	}
	return fmt.Errorf("doesn't exist")
}

// Delete an existing entity from local storage.
func (l *Local) Delete(ent Entity) error {
	for ii, entity := range l.Entities {
		if entity.ID() == ent.ID() {
			l.Entities[ii] = nil
			return nil
		}
	}
	return fmt.Errorf("doesn't exist")
}
