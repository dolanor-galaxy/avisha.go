// Package style contains stylised rendering for different widgets.
// Will attempt to follow material design language to fit in with materials package (gio).
package style

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
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
