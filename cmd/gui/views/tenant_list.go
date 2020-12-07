package views

import (
	"fmt"
	"sync"
	"unsafe"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
)

// Tenants displays a list of Tenants and provides controls for editing them.
type Tenants struct {
	nav.Route

	App *avisha.App
	Th  *style.Theme

	RegisterTenant widget.Clickable

	list   layout.List
	states States
	once   sync.Once
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
		t.Route.To(RouteTenantForm, nil)
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
	if err := t.App.DB.All(&tenants); err != nil {
		fmt.Printf("reading tenants: %s\n", err)
	}
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
