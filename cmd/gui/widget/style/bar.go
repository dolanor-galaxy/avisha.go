package style

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/util"
)

// TopBar renders a title and the provided actions.
type TopBar struct {
	*material.Theme
	Height unit.Value
}

func (bar TopBar) Layout(gtx C, title string, actions ...layout.Widget) D {
	gtx.Constraints.Max.Y = gtx.Px(bar.Height)
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return util.Rect{
				Color: bar.Color.Primary,
				Size:  gtx.Constraints.Max,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			items := []layout.FlexChild{
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(unit.Dp(10)).Layout(
						gtx,
						func(gtx C) D {
							title := material.Label(bar.Theme, unit.Dp(24), title)
							title.Color = bar.Theme.Color.InvText
							return title.Layout(gtx)
						})
				}),
				layout.Flexed(1, func(gtx C) D {
					return D{Size: gtx.Constraints.Min}
				}),
			}
			// @Todo: handle overflow (when actions don't fit the bar).
			// - Detect overflow (dim calcs probably)
			// - Render icon button
			// - Display a list of overflowed actions when clicked
			//
			// @Todo: auto centering of action content.
			// Atm insets are hard-coded, thus actions have to know the bar height
			// and if bar height changes then actions would need to change.
			for _, action := range actions {
				action := action
				items = append(items, layout.Rigid(action))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, items...)
		}),
	)
}
