package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

// File storage.
// Manually populate with type initializers in order to serialize/deserialize
// concrete types properly.
type File struct {
	// Path to file for storage.
	Path string
	// Runtime type information.
	Types map[string]Entity
	// Data to serialize.
	Data
	// format the json output.
	format bool
}

// Data to serialize.
type Data struct {
	// One bucket per type: "[type-name][entity-id]value".
	Buckets map[string]map[string]json.RawMessage
	// Order to iterate buckets.
	BucketOrder []string
	// Order to iterate bucket values.
	ValueOrder map[string][]string
}

// FileStorage allocates a File storage solution.
func FileStorage(path string) *File {
	f := &File{
		Path:  path,
		Types: make(map[string]Entity),
		Data: Data{
			Buckets:     make(map[string]map[string]json.RawMessage),
			BucketOrder: []string{},
			ValueOrder:  make(map[string][]string),
		},
	}
	return f
}

// With creates a bucket to contain the runtime type information required
// for serialization.
// Type must be a pointer.
func (f *File) With(prototype Entity) *File {
	name := reflect.ValueOf(prototype).Elem().Type().Name()
	if name == "" {
		panic(fmt.Sprintf("type name is empty for %T\n", prototype))
	}
	f.Types[name] = prototype
	f.Buckets[name] = make(map[string]json.RawMessage)
	f.BucketOrder = append(f.BucketOrder, name)
	return f
}

// Format json output.
func (f *File) Format(format bool) *File {
	f.format = format
	return f
}

// MustLoad loads data from file and panics on error.
func (f *File) MustLoad() *File {
	if err := f.load(); err != nil {
		panic(fmt.Errorf("loading file storage: %w", err))
	}
	return f
}

// Query the file storage.
func (f *File) Query(plist ...func(Entity) bool) (Entity, bool) {
	// Loop over all the buckets, every entity will be tested.
	for _, T := range f.BucketOrder {
		bucket := f.Buckets[T]
		for _, by := range bucket {
			// Deserialize each entity using the appropriate type.
			ent, ok := reflect.New(reflect.ValueOf(f.Types[T]).Elem().Type()).Interface().(Entity)
			if !ok {
				continue
			}
			// Test against predicates.
			if err := json.Unmarshal(by, ent); err == nil {
				if ok := Apply(ent, plist...); ok {
					return ent, ok
				}
			}
		}
	}
	return nil, false
}

// List query the file storage.
func (f *File) List(plist ...func(Entity) bool) []Entity {
	var found []Entity
	// Loop over all the buckets, every entity will be tested.
	for _, T := range f.BucketOrder {
		bucket := f.Buckets[T]
		for _, order := range f.ValueOrder[T] {
			by := bucket[order]
			// Deserialize each entity using the appropriate type.
			ent, ok := reflect.New(reflect.ValueOf(f.Types[T]).Elem().Type()).Interface().(Entity)
			if !ok {
				continue
			}
			// Test against predicates.
			if err := json.Unmarshal(by, ent); err == nil {
				if ok := Apply(ent, plist...); ok {
					found = append(found, ent)
				}
			}
		}
	}
	return found
}

// Create an entity
func (f *File) Create(ent Entity) error {
	return f.Save(ent)
}

// Update an entity
func (f *File) Update(ent Entity) error {
	return f.Save(ent)
}

// Save will update an entity or create a new on if it does not exist.
func (f *File) Save(ent Entity) error {
	for T := range f.Types {
		if T == reflect.TypeOf(ent).Name() {
			v, err := json.Marshal(ent)
			if err != nil {
				return fmt.Errorf("serializing entity: %w", err)
			}
			_, ok := f.Buckets[T][ent.ID()]
			f.Buckets[T][ent.ID()] = v
			if !ok {
				f.ValueOrder[T] = append(f.ValueOrder[T], ent.ID())
			}
			break
		}
	}
	return f.save()
}

// Delete an entity from storage.
func (f *File) Delete(ent Entity) error {
	for T := range f.Types {
		if T == reflect.TypeOf(ent).Name() {
			delete(f.Buckets[T], ent.ID())
			break
		}
	}
	return f.save()
}

func (f *File) save() error {
	var (
		indent = ""
	)

	os.Remove(f.Path) // Overwrite existing.

	handle, err := f.prepare()
	if err != nil {
		return err
	}
	defer handle.Close()

	if f.format {
		indent = "\t"
	}

	by, err := json.MarshalIndent(f.Data, "", indent)
	if err != nil {
		return fmt.Errorf("serializing buckets: %w", err)
	}

	if _, err := io.Copy(handle, bytes.NewBuffer(by)); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func (f *File) load() error {
	handle, err := f.prepare()
	if err != nil {
		return err
	}
	defer handle.Close()
	by, err := ioutil.ReadAll(handle)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	if len(by) == 0 {
		return nil
	}
	if err := json.Unmarshal(by, &f.Data); err != nil {
		return fmt.Errorf("deserializing buckets: %w", err)
	}
	return nil
}

func (f *File) prepare() (*os.File, error) {
	handle, err := os.OpenFile(f.Path, os.O_RDWR, 0777)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0777); err != nil {
			return nil, fmt.Errorf("preparing directories: %w", err)
		}
		handle, err = os.Create(f.Path)
		if err != nil {
			return nil, fmt.Errorf("creating storage file: %w", err)
		}
	}
	return handle, nil
}
