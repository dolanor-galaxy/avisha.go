package views

import (
	"image"
	"log"
	"sync"
	"unsafe"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/icons"
	"github.com/jackmordaunt/avisha.go/cmd/gui/nav"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
)

// @Todo: Cap width for list items for desktop view and pack into columns?
// @Todo: Add search and filters
type LeaseList struct {
	nav.Route
	App    *avisha.App
	Th     *style.Theme
	list   layout.List
	states States
	once   sync.Once

	CreateLease widget.Clickable
}

func (l *LeaseList) Title() string {
	return "Leases"
}

func (l *LeaseList) Context() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return material.IconButton(
				l.Th.Primary(),
				&l.CreateLease,
				icons.Add,
			).Layout(gtx)
		},
	}
}

func (l *LeaseList) Update(gtx C) {
	for _, state := range l.states.List() {
		if state.Item.Clicked() {
			l.Route.To(RouteLeasePage, (*avisha.Lease)(state.Data))
		}
	}
	if l.CreateLease.Clicked() {
		l.Route.To(RouteLeasePage, nil)
	}
}

func (l *LeaseList) Layout(gtx C) D {
	l.once.Do(func() {
		l.list.Axis = layout.Vertical
		l.list.ScrollToEnd = false
	})
	l.Update(gtx)
	l.states.Begin()
	var (
		leases []*avisha.Lease
	)
	if err := l.App.All(&leases); err != nil {
		log.Printf("loading leases: %v", err)
	}
	return l.list.Layout(gtx, len(leases), func(gtx C, index int) D {
		var (
			lease  = leases[index]
			state  = l.states.Next(unsafe.Pointer(lease))
			active = false
			tenant avisha.Tenant
			site   avisha.Site
		)
		// @Todo: handle data loading errors by
		// - Displaying a message to the user that something went wrong.
		// - Logging error (it's a system bug, not a user error) to service.
		// - Generate bug report?
		// Is it worth trying to abstract this data loading stuff?
		if err := l.App.One("ID", lease.Tenant, &tenant); err != nil {
			log.Printf("lease list: %v", err)
		}
		if err := l.App.One("ID", lease.Site, &site); err != nil {
			log.Printf("lease list: %v", err)
		}
		return style.ListItem(
			gtx,
			l.Th.Dark(),
			&state.Item,
			&state.Hover,
			active,
			func(gtx C) D {
				return style.Card{
					Content: []layout.Widget{
						func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(
								gtx,
								layout.Rigid(func(gtx C) D {
									return style.TenantLabel(l.Th.Dark(), tenant).Layout(gtx)
								}),
								layout.Flexed(1.0, func(gtx C) D {
									return D{Size: image.Point{
										X: gtx.Px(unit.Dp(5)),
										Y: 0,
									}}
								}),
								layout.Rigid(func(gtx C) D {
									return material.Label(
										l.Th.Dark(),
										unit.Dp(20),
										site.Number,
									).Layout(gtx)
								}),
							)
						},
						func(gtx C) D {
							u := lease.Services["utilities"]
							return style.ServiceLabel(l.Th, "Utilities", u.Balance()).Layout(gtx)
						},
						func(gtx C) D {
							u := lease.Services["rent"]
							return style.ServiceLabel(l.Th, "Rent", u.Balance()).Layout(gtx)
						},
						func(gtx C) D {
							lb := material.Label(
								l.Th.Muted(),
								unit.Dp(15),
								lease.Term.String())
							// lb.Color = l.Th.Muted().Fg
							return lb.Layout(gtx)
						},
					},
				}.Layout(gtx, l.Th.Dark())
			})
	})
}
