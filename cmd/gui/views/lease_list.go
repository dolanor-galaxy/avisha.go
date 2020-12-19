package views

import (
	"fmt"
	"image"
	"log"
	"sort"
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
	// @Improve
	// Using this composite type so we can order by site number, since
	// we can't do alphabetic ordering of integer IDs.
	// Not sure about performance hit this takes, it's at least O(2n).
	type Lease struct {
		Lease  avisha.Lease
		Site   avisha.Site
		Tenant avisha.Tenant
	}
	var (
		list = make([]Lease, len(leases))
	)
	for ii := range leases {
		list[ii].Lease = *leases[ii]
		if err := l.App.One("ID", leases[ii].Tenant, &list[ii].Tenant); err != nil {
			log.Printf("lease list: %v", err)
		}
		if err := l.App.One("ID", leases[ii].Site, &list[ii].Site); err != nil {
			log.Printf("lease list: %v", err)
		}
	}
	sort.Slice(list, func(ii, jj int) bool {
		return list[ii].Site.Number < list[jj].Site.Number
	})
	return l.list.Layout(gtx, len(list), func(gtx C, index int) D {
		var (
			lease  = &list[index].Lease
			tenant = &list[index].Tenant
			site   = &list[index].Site
			state  = l.states.Next(unsafe.Pointer(lease))
			active = false
		)
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
									return material.Label(
										l.Th.Dark(),
										unit.Dp(20),
										fmt.Sprintf("Site %v", site.Number),
									).Layout(gtx)
								}),
								layout.Flexed(1.0, func(gtx C) D {
									return D{Size: image.Point{
										X: gtx.Px(unit.Dp(5)),
										Y: 0,
									}}
								}),
								layout.Rigid(func(gtx C) D {
									return style.TenantLabel(l.Th.Dark(), *tenant).Layout(gtx)
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
