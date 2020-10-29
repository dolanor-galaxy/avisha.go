package views

import (
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
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/theme"
	"github.com/jackmordaunt/avisha-fn/storage"
)

// Tenants displays a list of Tenants and provides controls for editing them.
type Tenants struct {
	nav.Route
	App    *avisha.App
	Th     *theme.Theme
	list   layout.List
	states States
	once   sync.Once

	RegisterTenant widget.Clickable
}

func (t *Tenants) Title() string {
	return "Tenants"
}

func (t *Tenants) Context() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return material.IconButton(
				t.Th.Primary(),
				&t.RegisterTenant,
				icons.Add,
			).Layout(gtx)
		},
	}
}

func (t *Tenants) Update(gtx C) {
	for _, state := range t.states.List() {
		if state.Item.Clicked() {
			t.Route.To(RouteTenantForm, (*avisha.Tenant)(state.Data))
		}
	}
	if t.RegisterTenant.Clicked() {
		t.Route.To(RouteTenantForm, &avisha.Tenant{})
	}
}

func (t *Tenants) Layout(gtx C) D {
	t.once.Do(func() {
		t.list.Axis = layout.Vertical
		t.list.ScrollToEnd = false
	})
	t.Update(gtx)
	t.states.Begin()
	var (
		tenants []*avisha.Tenant
	)
	t.App.List(func(ent storage.Entity) bool {
		t, ok := ent.(*avisha.Tenant)
		if ok {
			tenants = append(tenants, t)
		}
		return ok
	})
	return t.list.Layout(gtx, len(tenants), func(gtx C, index int) D {
		var (
			tenant = tenants[index]
			state  = t.states.Next(unsafe.Pointer(tenant))
			active = false
		)
		return style.ListItem(
			gtx,
			t.Th.Primary(),
			&state.Item,
			&state.Hover,
			active,
			func(gtx C) D {
				return material.Label(
					t.Th.Primary(),
					unit.Dp(20),
					tenant.Name,
				).Layout(gtx)
			})
	})
}
