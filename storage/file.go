package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Type wraps a concrete type, letting us allocate new objects and test for
// type equality.
//
// Note: Since I see no good way of specifying types from callsite,
// I'm resorting to a hacky work-around: use caller-defined "Type" objects
// that can:
//
// 	1. Allocate new instances of the concrete type.
// 	2. Confirm if a given concrete object matches it's type.
//
type Type interface {
	New() interface{}
	Is(interface{}) bool
}

// File storage.
// Manually populate with type initializers in order to serialize/deserialize
// concrete types properly.
type File struct {
	// Path to file for storage.
	Path string
	// Map of type name to it's type initializer.
	Types map[string]Type
	// Buckets groups each type of object into it's own bucket
	// so that we know what concrete type to deserialize into.
	Buckets map[string][]json.RawMessage
}

// Query the file storage.
func (f *File) Query(plist []Predicate) (interface{}, bool) {
	// TODO: Handle file errors during initialization.
	by, _ := ioutil.ReadFile(f.Path)
	_ = json.Unmarshal(by, &f.Buckets)
	// Loop over all the buckets, every entity will be tested.
	for kind, b := range f.Buckets {
		for _, msg := range b {
			// Deserialize each entity using the appropriate type.
			ent := f.Types[kind].New()
			if err := json.Unmarshal(msg, ent); err == nil {
				if ok := Predicates(plist).Apply(ent); ok {
					return ent, ok
				}
			}
		}
	}
	return nil, false
}

// Save to file.
func (f *File) Save(ent interface{}) error {
	for kind, t := range f.Types {
		if t.Is(ent) {
			v, err := json.Marshal(ent)
			if err != nil {
				return fmt.Errorf("serializing entity: %w", err)
			}
			f.Buckets[kind] = append(f.Buckets[kind], v)
			break
		}
	}
	by, err := json.Marshal(f.Buckets)
	if err != nil {
		return fmt.Errorf("serializing buckets: %w", err)
	}
	if err := ioutil.WriteFile(f.Path, by, os.ModePerm); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	return nil
}
