package views

import (
	"unsafe"

	"gioui.org/layout"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Route = string

const (
	RouteLease      Route = "lease"
	RouteTenants    Route = "tenants"
	RouteLeaseForm  Route = "lease-form"
	RouteTenantForm Route = "tenant-form"
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

func (s *States) List() []State {
	return s.list[:s.current]
}
