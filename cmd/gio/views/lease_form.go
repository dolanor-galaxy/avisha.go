package views

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
)

type LeaseForm struct {
	*avisha.App
	*material.Theme
	lease *avisha.Lease
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

func (l *LeaseForm) Layout(gtx Ctx) Dims {
	return Dims{}
}
