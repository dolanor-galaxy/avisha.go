// Package style contains stylised rendering for different widgets.
// Will attempt to follow material design language to fit in with materials package (gio).
package style

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// ListItem renders a list item.
func ListItem(
	gtx C,
	th *material.Theme,
	state *widget.Clickable,
	hover *widget.Hoverable,
	active bool,
	w layout.Widget,
) D {
	return layout.Stack{}.Layout(
		gtx,
		layout.Expanded(func(gtx C) D {
			if active {
				util.Rect{
					Color: WithAlpha(th.Color.Hint, 150),
					Size:  gtx.Constraints.Max,
				}.Layout(gtx)
			} else if hover.Hovered() {
				util.Rect{
					Color: WithAlpha(th.Color.Hint, 38),
					Size:  gtx.Constraints.Max,
				}.Layout(gtx)
			}
			return D{}
		}),
		layout.Stacked(func(gtx C) D {
			return material.Clickable(gtx, state, func(gtx C) D {
				return layout.UniformInset(unit.Dp(10)).Layout(
					gtx,
					func(gtx C) D {
						dims := w(gtx)
						dims.Size.X = gtx.Constraints.Max.X
						return dims
					},
				)
			})
		}),
		layout.Expanded(func(gtx C) D {
			return hover.Layout(gtx)
		}),
	)
}

func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

// Dialog renders an input with ok / cancel actions.
type Dialog struct {
	Context string

	Input  materials.TextField
	Ok     widget.Clickable
	Cancel widget.Clickable
}

func (d *Dialog) Layout(gtx C, th *material.Theme, title string) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return d.Input.Layout(gtx, th, title)
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{Y: gtx.Px(unit.Dp(10))}}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(
				gtx,
				layout.Flexed(1, func(gtx C) D {
					return D{Size: gtx.Constraints.Min}
				}),
				layout.Rigid(func(gtx C) D {
					btn := material.Button(th, &d.Cancel, "Cancel")
					btn.Color = btn.Background
					btn.Background = WithAlpha(btn.Background, 0)
					return btn.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
				}),
				layout.Rigid(func(gtx C) D {
					return material.Button(th, &d.Ok, "Ok").Layout(gtx)
				}),
			)
		}),
	)
}
