package views

import (
	"unsafe"

	"gioui.org/layout"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Route = string

const (
	RouteLease      Route = "lease"
	RouteTenants    Route = "tenants"
	RouteSites      Route = "sites"
	RouteLeasePage  Route = "lease-page"
	RouteTenantForm Route = "tenant-form"
	RouteSiteForm   Route = "site-form"
)

// States maintains list-item state, between frame updates.
// Allocates memory "as needed", alternative to pre-allocating
// slice of states.
// Arbitrary data can be stored behind a raw pointer.
type States struct {
	current int
	list    []State
}

type State struct {
	Data  unsafe.Pointer
	Item  widget.Clickable
	Hover widget.Hoverable
}

func (s *States) Begin() {
	s.current = 0
}

func (s *States) Next(data unsafe.Pointer) *State {
	defer func() { s.current += 1 }()
	if s.current > len(s.list)-1 {
		s.list = append(s.list, State{})
	}
	state := &s.list[s.current]
	state.Data = data
	return state
}

// List returns a view over the active states in memory.
func (s *States) List() []*State {
	list := make([]*State, s.current)
	view := s.list[:s.current]
	for ii := range view {
		list[ii] = &view[ii]
	}
	return list
}
