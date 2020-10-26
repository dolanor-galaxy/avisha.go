package views

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

type LeaseForm struct {
	*avisha.App
	*material.Theme
	lease *avisha.Lease

	Tenant materials.TextField
	Site   materials.TextField
	Date   materials.TextField
	Term   materials.TextField
	// Start Date
	// Duration Int
	Submit widget.Clickable
	layout.List
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
		fmt.Printf("submitted: tenant %q, site %q\n", l.Tenant.Text(), l.Site.Text())
	}
}

func (l *LeaseForm) Layout(gtx Ctx) Dims {
	l.Update(gtx)
	l.List.Axis = layout.Vertical
	return layout.UniformInset(unit.Dp(10)).Layout(
		gtx,
		func(gtx Ctx) Dims {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx Ctx) Dims {
					return l.Tenant.Layout(gtx, l.Theme, "Tenant")
				}),
				layout.Rigid(func(gtx Ctx) Dims {
					return l.Date.Layout(gtx, l.Theme, "Date")
				}),
				layout.Rigid(func(gtx Ctx) Dims {
					return l.Site.Layout(gtx, l.Theme, "Site")
				}),
			)
		},
	)
}
