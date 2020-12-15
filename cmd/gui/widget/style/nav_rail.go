package style

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
)

// NavRail implements material.io navigation rail.
type NavRail struct {
	Th           *Theme
	Destinations []Destination
	Width        unit.Value
	layout.List
}

// Destination is a nav item in a rail.
// It can have an icon, label and is clickable.
type Destination struct {
	widget.Clickable
	Active bool
	Route  string
	Label  string
	Icon   *widget.Icon
}

func (r *NavRail) Layout(gtx C) D {
	r.List.Axis = layout.Vertical
	width := gtx.Px(r.Width)
	cs := &gtx.Constraints
	cs.Max.X = width
	cs.Min.X = width
	cs.Min.Y = gtx.Constraints.Max.Y
	// Draw vertical line.
	// Not sure on the best approach.
	// This code just clips a rectangle.
	stack := op.Push(gtx.Ops)
	clip.Rect{
		Max: image.Point{
			Y: cs.Max.Y,
			X: width,
		},
		Min: image.Point{
			Y: 0,
			X: width - gtx.Px(unit.Dp(2)),
		},
	}.Add(gtx.Ops)
	clip.Rect{Max: cs.Max}.Add(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{A: 100}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Pop()
	cs.Max.X -= gtx.Px(unit.Dp(1))
	cs.Min.X -= gtx.Px(unit.Dp(1))
	if len(r.Destinations) == 0 {
		return D{Size: image.Point{X: width, Y: cs.Max.Y}}
	}
	return r.List.Layout(gtx, len(r.Destinations), func(gtx C, ii int) D {
		var (
			item = &r.Destinations[ii]
		)
		return material.Clickable(gtx, &item.Clickable, func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(10),
				Right:  unit.Dp(10),
			}.Layout(
				gtx,
				func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							if item.Icon == nil {
								return D{}
							}
							item.Icon.Color = r.Th.Fg
							if item.Active {
								item.Icon.Color = r.Th.Fg
							}
							return item.Icon.Layout(gtx, unit.Dp(25))
						}),
						layout.Rigid(func(gtx C) D {
							l := material.Label(r.Th.Dark(), unit.Dp(16), item.Label)
							l.Alignment = text.Middle
							l.Color = r.Th.Fg
							if item.Active {
								l.Color = r.Th.Fg
							}
							return l.Layout(gtx)
						}),
					)
				})
		})
	})
}
