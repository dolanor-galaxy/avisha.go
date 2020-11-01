// Package style contains stylised rendering for different widgets.
// Will attempt to follow material design language to fit in with materials package (gio).
package style

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
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
					Color: materials.AlphaMultiply(th.Color.Hint, 150),
					Size:  gtx.Constraints.Max,
				}.Layout(gtx)
			} else if hover.Hovered() {
				util.Rect{
					Color: materials.AlphaMultiply(th.Color.Hint, 38),
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
