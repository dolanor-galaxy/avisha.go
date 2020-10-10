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
	Radii float32
}

func (r Rect) Layout(gtx C) D {
	return DrawRect(gtx, r.Color, r.Size, r.Radii)
}

// DrawRect creates a rectangle of the provided background color with
// Dimensions specified by size and a corner radius (on all corners)
// specified by radii.
func DrawRect(gtx C, background color.RGBA, size image.Point, radii float32) D {
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

// Widget can render itself.
type Widget = func(gtx C) Dimensions

// Dimensions of a rendered widget.
// Includes a macro for overlays.
type Dimensions struct {
	layout.Dimensions
	Overlay op.CallOp
}

// Flex implements a flex layout.
// Any widget that implements Overlayable will be overlayed on top.
// Flex wraps `layout.Flex` for underlying flex behaviour.
type Flex struct {
	layout.Flex
}

type FlexChild struct {
	Widget Widget
	Flex   bool
	Weight float32
}

func Rigid(w Widget) FlexChild {
	return FlexChild{
		Widget: w,
	}
}

func Flexed(weight float32, w Widget) FlexChild {
	return FlexChild{
		Widget: w,
		Weight: weight,
		Flex:   true,
	}
}

func (f Flex) Layout(gtx C, children ...FlexChild) D {
	var dimlist = make([]Dimensions, len(children))
	return layout.Stack{}.Layout(
		gtx,
		layout.Stacked(func(gtx C) D {
			var flex []layout.FlexChild
			for ii, child := range children {
				ii := ii
				child := child
				if child.Flex {
					flex = append(flex, layout.Flexed(child.Weight, func(gtx C) D {
						dims := child.Widget(gtx)
						dimlist[ii] = dims
						return dims.Dimensions
					}))
				} else {
					flex = append(flex, layout.Rigid(func(gtx C) D {
						dims := child.Widget(gtx)
						dimlist[ii] = dims
						return dims.Dimensions
					}))
				}
			}
			return f.Flex.Layout(gtx, flex...)
		}),
		layout.Expanded(func(gtx C) D {
			for ii := len(children) - 1; ii >= 0; ii-- {
				var (
					offset = image.Point{}
					dims   = dimlist[ii]
				)
				for jj := 0; jj <= ii; jj++ {
					offset = offset.Add(dimlist[jj].Size)
				}
				// TODO: Support overlaying all 4 sides.
				layout.Inset{
					Top: unit.Px(float32(offset.Y)),
				}.Layout(
					gtx,
					func(gtx C) D {
						dims.Overlay.Add(gtx.Ops)
						return D{}
					})
			}
			return D{}
		}),
	)
}
