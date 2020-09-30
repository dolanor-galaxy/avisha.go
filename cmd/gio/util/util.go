package util

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type (
	Ctx  = layout.Context
	Dims = layout.Dimensions
)

type Rect struct {
	Color color.RGBA
	Size  image.Point
	Radii float32
}

func (r Rect) Layout(gtx Ctx) Dims {
	return DrawRect(gtx, r.Color, r.Size, r.Radii)
}

// DrawRect creates a rectangle of the provided background color with
// Dimensions specified by size and a corner radius (on all corners)
// specified by radii.
func DrawRect(gtx Ctx, background color.RGBA, size image.Point, radii float32) Dims {
	stack := op.Push(gtx.Ops)
	{
		paint.ColorOp{
			Color: background,
		}.Add(gtx.Ops)
		bounds := f32.Rectangle{
			Max: layout.FPt(size),
		}
		if radii != 0 {
			clip.RRect{
				Rect: bounds,
				NW:   radii,
				NE:   radii,
				SE:   radii,
				SW:   radii,
			}.Add(gtx.Ops)
		}
		paint.PaintOp{
			Rect: bounds,
		}.Add(gtx.Ops)
	}
	stack.Pop()
	return layout.Dimensions{Size: size}
}

// TopBar renders a title and the provided actions.
type TopBar struct {
	*material.Theme
}

func (bar TopBar) Layout(gtx Ctx, title string, actions ...layout.Widget) Dims {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx Ctx) Dims {
			return Rect{
				Color: color.RGBA{B: 100, R: 50, A: 255},
				Size:  gtx.Constraints.Max,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx Ctx) Dims {
			items := []layout.FlexChild{
				layout.Flexed(float32(len(actions)+1), func(gtx Ctx) Dims {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx Ctx) Dims {
						title := material.Label(bar.Theme, unit.Dp(24), title)
						title.Color = bar.Theme.Color.InvText
						return title.Layout(gtx)
					})
				}),
			}
			for _, action := range actions {
				action := action
				items = append(items, layout.Flexed(1, func(gtx Ctx) Dims {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, action)
				}))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, items...)
		}),
	)
}
