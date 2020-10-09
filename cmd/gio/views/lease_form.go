package views

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
)

type LeaseForm struct {
	*avisha.App
	*material.Theme
	lease *avisha.Lease

	Tenant style.Select
	Site   style.Select
	Foo    style.Select
	Date   style.TextField
	Term   style.TextField
	// Start Date
	// Duration Int
	Submit widget.Clickable
}

func (l *LeaseForm) ReRoute() (string, bool) {
	return "", false
}

func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok {
		l.lease = lease
	}
}

func (l *LeaseForm) Context() (list []layout.Widget) {
	if l.lease != nil {
		list = append(list, func(gtx Ctx) Dims {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx Ctx) Dims {
					label := material.Label(l.Theme, unit.Dp(24), l.lease.ID())
					label.Alignment = text.Middle
					label.Color = l.Theme.Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (l *LeaseForm) Update(gtx Ctx) {
	if l.Submit.Clicked() {
		// grab data and submit to app;
		fmt.Printf("submitted: tenant %q, site %q\n", l.Tenant.Value, l.Site.Value)
	}
}

func (l *LeaseForm) Layout(gtx Ctx) Dims {
	l.Update(gtx)
	return layout.UniformInset(unit.Dp(10)).Layout(
		gtx,
		func(gtx Ctx) Dims {
			return util.Flex{
				Flex: layout.Flex{Axis: layout.Vertical},
			}.Layout(
				gtx,
				func(gtx Ctx) util.OverlayChild {
					gtx.Constraints.Max.X = gtx.Px(unit.Dp(80))
					return l.Tenant.Layout(gtx, l.Theme, "Tenant", []string{"one", "two", "three"})
				}(gtx),
				l.Site.Layout(gtx, l.Theme, "Site", []string{"one", "two", "three"}),
				l.Foo.Layout(gtx, l.Theme, "Foo", []string{"one", "two", "three"}),
			)
		},
	)
}
