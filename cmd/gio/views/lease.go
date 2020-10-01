package views

import (
	"fmt"
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
	"github.com/jackmordaunt/avisha-fn/storage"
)

type Lease struct {
	*avisha.App
	*material.Theme

	list   layout.List
	states *States

	route    string
	selected *avisha.Lease
	once     sync.Once
}

func (l *Lease) ReRoute() (string, interface{}) {
	if l.route != "" && l.selected != nil {
		defer func() { l.route = "" }()
		return l.route, l.selected
	}
	return "", nil
}

func (l *Lease) Context() []layout.Widget {
	return nil
}

func (l *Lease) Update(gtx Ctx) {
	for _, state := range l.states.List() {
		if state.Item.Clicked() {
			fmt.Printf("navigating to LeaseForm for %s\n", state.ID())
			l.route = "LeaseForm"
			l.selected = state.Lease
		}
		// if state.Hover.Hovered() {
		// 	// fmt.Printf("%s is hovered\n", state.ID())
		// }
	}
}

func (l *Lease) Layout(gtx Ctx) Dims {
	l.once.Do(func() {
		l.list.Axis = layout.Vertical
		l.list.ScrollToEnd = false
		l.states = &States{}
	})
	l.Update(gtx)
	l.states.Begin()
	var (
		leases []*avisha.Lease
	)
	l.App.List(func(ent storage.Entity) bool {
		l, ok := ent.(*avisha.Lease)
		if ok {
			leases = append(leases, l)
		}
		return ok
	})
	return l.list.Layout(gtx, len(leases), func(gtx Ctx, index int) Dims {
		var (
			lease  = leases[index]
			state  = l.states.Next(lease)
			active = false
		)
		return style.ListItem(gtx, l.Theme, &state.Item, &state.Hover, active, func(gtx Ctx) Dims {
			return material.Label(l.Theme, unit.Dp(20), fmt.Sprintf("%s - %s: %+v", lease.Site, lease.Tenant, lease.Term)).Layout(gtx)
		})
	})
}

// States tracks state per list-item, between frame updates.
// Allocates memory "as needed", alternative to pre-allocating
// slice of states.
type States struct {
	current int
	list    []State
}

type State struct {
	*avisha.Lease
	Item  widget.Clickable
	Hover widget.Hoverable
}

func (s *States) Begin() {
	s.current = 0
}

func (s *States) Next(lease *avisha.Lease) *State {
	defer func() { s.current += 1 }()
	if s.current > len(s.list)-1 {
		s.list = append(s.list, State{})
	}
	state := &s.list[s.current]
	state.Lease = lease
	return state
}

func (s *States) List() []State {
	return s.list[:s.current]
}
