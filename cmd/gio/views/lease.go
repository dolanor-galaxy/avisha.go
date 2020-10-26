package views

import (
	"fmt"
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
	"github.com/jackmordaunt/avisha-fn/storage"
)

type Lease struct {
	nav.Route
	*avisha.App
	*material.Theme

	CreateLease widget.Clickable

	list   layout.List
	states *States

	once sync.Once
}

func (l *Lease) Title() string {
	return "Leases"
}

func (l *Lease) Receive(v interface{}) {
	l.states = &States{}
}

func (l *Lease) Context() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return material.IconButton(l.Theme, &l.CreateLease, icons.Add).Layout(gtx)
		},
	}
}

func (l *Lease) Update(gtx C) {
	for _, state := range l.states.List() {
		for state.Item.Clicked() {
			l.Route.To(RouteLeaseForm, state.Lease)
		}
	}
	if l.CreateLease.Clicked() {
		l.Route.To(RouteLeaseForm, &avisha.Lease{})
	}
}

func (l *Lease) Layout(gtx C) D {
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
	return l.list.Layout(gtx, len(leases), func(gtx C, index int) D {
		var (
			lease  = leases[index]
			state  = l.states.Next(lease)
			active = false
		)
		return style.ListItem(gtx, l.Theme, &state.Item, &state.Hover, active, func(gtx C) D {
			return material.Label(
				l.Theme,
				unit.Dp(20),
				fmt.Sprintf("%s - %s: %+v", lease.Site, lease.Tenant, lease.Term),
			).Layout(gtx)
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
