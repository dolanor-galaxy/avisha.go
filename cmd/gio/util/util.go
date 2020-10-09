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

// Widget types can render themselves.
type Widget interface {
	Layout(gtx Ctx) Dims
}

// Overlayable types can render an overlay.
type Overlayable interface {
	Overlay(gtx Ctx) Dims
}

// WidgetFunc implements Widget.
type WidgetFunc func(gtx Ctx) Dims

func (w WidgetFunc) Layout(gtx Ctx) Dims {
	return w(gtx)
}

// OverlayChild allows for convenient implementation of Widget and Overlayable
// with closures.
type OverlayChild struct {
	Content   layout.Widget
	Overlayed layout.Widget
}

func (c OverlayChild) Layout(gtx Ctx) Dims {
	return c.Content(gtx)
}

func (c OverlayChild) Overlay(gtx Ctx) Dims {
	if c.Overlayed == nil {
		return Dims{}
	}
	return c.Overlayed(gtx)
}

// Overlayed constructs an overlay widget from two closures.
func Overlayed(w layout.Widget, o layout.Widget) OverlayChild {
	return OverlayChild{
		Content:   w,
		Overlayed: o,
	}
}

// Flex implements a flex layout.
// Any widget that implements Overlayable will be overlayed on top.
// Flex wraps `layout.Flex` for underlying flex behaviour.
type Flex struct {
	layout.Flex
}

func (f Flex) Layout(gtx Ctx, children ...Widget) Dims {
	var dimlist = make([]Dims, len(children))
	return layout.Stack{}.Layout(
		gtx,
		layout.Stacked(func(gtx Ctx) Dims {
			var flex []layout.FlexChild
			for ii, child := range children {
				ii := ii
				child := child
				flex = append(flex, layout.Rigid(func(gtx Ctx) Dims {
					dims := child.Layout(gtx)
					dimlist[ii] = dims
					return dims
				}))
			}
			return f.Flex.Layout(gtx, flex...)
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			for ii := len(children) - 1; ii >= 0; ii-- {
				var (
					ii     = ii
					offset = image.Point{}
					child  = children[ii]
				)
				overlay, ok := child.(Overlayable)
				if !ok {
					continue
				}
				for jj := 0; jj <= ii; jj++ {
					offset = offset.Add(dimlist[jj].Size)
				}
				// TODO: Support overlaying all 4 sides.
				layout.Inset{
					Top: unit.Px(float32(offset.Y)),
				}.Layout(
					gtx,
					func(gtx Ctx) Dims {
						return overlay.Overlay(gtx)
					})
			}
			return Dims{}
		}),
	)
}
