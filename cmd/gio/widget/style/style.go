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
	Ctx  = layout.Context
	Dims = layout.Dimensions
)

// ListItem renders a list item.
func ListItem(gtx Ctx, th *material.Theme, state *widget.Clickable, hover *widget.Hoverable, active bool, w layout.Widget) Dims {
	return layout.Stack{}.Layout(
		gtx,
		layout.Stacked(func(gtx Ctx) Dims {
			if active {
				util.Rect{
					Color: materials.AlphaMultiply(th.Color.Hint, 150),
					Size:  gtx.Constraints.Max,
				}.Layout(gtx)
			} else if hover.Hovered() {
				util.Rect{
					Color: materials.AlphaMultiply(th.Color.Hint, 100),
					Size:  gtx.Constraints.Max,
				}.Layout(gtx)
			}
			return Dims{}
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			return material.Clickable(gtx, state, func(gtx Ctx) Dims {
				return layout.UniformInset(unit.Dp(10)).Layout(
					gtx,
					func(gtx Ctx) Dims {
						dims := w(gtx)
						dims.Size.X = gtx.Constraints.Max.X
						return dims
					},
				)
			})
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			return hover.Layout(gtx)
		}),
	)
}
