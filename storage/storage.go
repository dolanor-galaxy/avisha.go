package storage

// Storer provides persistence for entities.
type Storer interface {
	Query(filters ...func(interface{}) bool) (interface{}, bool)
	List(filters ...func(interface{}) bool) []interface{}
	Save(ent interface{}) error
}

// Apply all filters to ent.
func Apply(ent interface{}, filters ...func(interface{}) bool) bool {
	for _, filter := range filters {
		if !filter(ent) {
			return false
		}
	}
	return true
}

// // Tenant storage.
// type Tenant interface {
// 	Tenants(filters ...func(*avisha.Tenant) bool) []avisha.Tenant
// }

// // Site storage.
// type Site interface {
// 	Sites(filters ...func(*avisha.Site) bool) []avisha.Site
// }

// // Lease storage.
// type Lease interface {
// 	Leases(filters ...func(*avisha.Lease) bool) []avisha.Lease
// }
