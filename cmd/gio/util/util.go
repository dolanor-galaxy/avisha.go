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
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Rect struct {
	Color color.RGBA
	Size  image.Point
	Radii unit.Value
}

func (r Rect) Layout(gtx C) D {
	return DrawRect(gtx, r.Color, r.Size, r.Radii)
}

// DrawRect creates a rectangle of the provided background color with
// Dimensions specified by size and a corner radius (on all corners)
// specified by radii.
func DrawRect(gtx C, background color.RGBA, size image.Point, radii unit.Value) D {
	defer op.Push(gtx.Ops).Pop()
	rr := float32(gtx.Px(radii))
	paint.ColorOp{
		Color: background,
	}.Add(gtx.Ops)
	bounds := f32.Rectangle{
		Max: layout.FPt(size),
	}
	if rr != 0 {
		clip.RRect{
			Rect: bounds,
			NW:   rr,
			NE:   rr,
			SE:   rr,
			SW:   rr,
		}.Add(gtx.Ops)
	}
	paint.PaintOp{
		Rect: bounds,
	}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}
