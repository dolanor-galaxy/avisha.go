package views

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
)

type LeaseForm struct {
	*avisha.App
	*material.Theme
}

func (l *LeaseForm) ReRoute() (string, bool) {
	return "", false
}

func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok {
		fmt.Printf("received lease: %v\n", lease)
	}
}

func (l *LeaseForm) Context() []layout.Widget {
	return []layout.Widget{
		func(gtx Ctx) Dims {
			label := material.Label(l.Theme, l.Theme.TextSize, "Context")
			label.Alignment = text.Middle
			label.Color = l.Theme.Color.InvText
			return label.Layout(gtx)
		},
		// func(gtx Ctx) Dims {
		// 	return util.Rect{Size: gtx.Constraints.Max, Color: color.RGBA{R: 255, A: 255}}.Layout(gtx)
		// },
	}
}

func (l *LeaseForm) Layout(gtx Ctx) Dims {
	return Dims{}
}
