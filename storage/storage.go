package storage

// Entity is a unique object.
type Entity interface {
	ID() string
}

// Storage is any object that can store and query entities.
type Storage interface {
	Storer
	Queryer
}

// Storer provides persistence for entities.
type Storer interface {
	// Create a new, unique entity.
	Create(ent Entity) error
	// Update an existing entity.
	Update(ent Entity) error
	// Delete an existing entity.
	Delete(ent Entity) error
}

// Queryer is an object that can be queried for entities.
type Queryer interface {
	// Query for a single entity that satisfies the filters.
	Query(filters ...func(Entity) bool) (Entity, bool)
	// List all entities that satisfy the filters.
	List(filters ...func(Entity) bool) []Entity
}

// Apply all filters to the entity.
func Apply(entity Entity, filters ...func(Entity) bool) bool {
	for _, filter := range filters {
		if !filter(entity) {
			return false
		}
	}
	return true
}
