package views

import (
	"fmt"
	"sync"
	"unsafe"

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

type Sites struct {
	nav.Route
	*avisha.App
	*material.Theme

	RegisterSite widget.Clickable

	list   layout.List
	states *States
	once   sync.Once
}

func (s *Sites) Title() string {
	return "Leases"
}

func (s *Sites) Receive(v interface{}) {
	s.states = &States{}
}

func (s *Sites) Context() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return material.IconButton(s.Theme, &s.RegisterSite, icons.Add).Layout(gtx)
		},
	}
}

func (s *Sites) Update(gtx C) {
	for _, state := range s.states.List() {
		for state.Item.Clicked() {
			s.Route.To(RouteSiteForm, (*avisha.Site)(state.Data))
		}
	}
	if s.RegisterSite.Clicked() {
		s.Route.To(RouteSiteForm, &avisha.Site{})
	}
}

func (s *Sites) Layout(gtx C) D {
	s.once.Do(func() {
		s.list.Axis = layout.Vertical
		s.list.ScrollToEnd = false
		s.states = &States{}
	})
	s.Update(gtx)
	s.states.Begin()
	var (
		sites []*avisha.Site
	)
	s.App.List(func(ent storage.Entity) bool {
		site, ok := ent.(*avisha.Site)
		if ok {
			sites = append(sites, site)
		}
		return ok
	})
	return s.list.Layout(gtx, len(sites), func(gtx C, index int) D {
		var (
			site   = sites[index]
			state  = s.states.Next(unsafe.Pointer(site))
			active = false
		)
		return style.ListItem(gtx, s.Theme, &state.Item, &state.Hover, active, func(gtx C) D {
			return material.Label(
				s.Theme,
				unit.Dp(20),
				fmt.Sprintf("%s", site.Number),
			).Layout(gtx)
		})
	})
}
