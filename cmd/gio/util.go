package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
