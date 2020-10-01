package style

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/util"
)

// TopBar renders a title and the provided actions.
type TopBar struct {
	*material.Theme
}

func (bar TopBar) Layout(gtx Ctx, title string, actions ...layout.Widget) Dims {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx Ctx) Dims {
			return util.Rect{
				Color: bar.Color.Primary,
				Size:  gtx.Constraints.Max,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx Ctx) Dims {
			items := []layout.FlexChild{
				layout.Rigid(func(gtx Ctx) Dims {
					return layout.UniformInset(unit.Dp(10)).Layout(
						gtx,
						func(gtx Ctx) Dims {
							title := material.Label(bar.Theme, unit.Dp(24), title)
							title.Color = bar.Theme.Color.InvText
							return title.Layout(gtx)
						})
				}),
			}
			// TODO: handle overflow (when actions don't fit the bar).
			// - Detect overflow (dim calcs probably)
			// - Render icon button
			// - Display a list of overflowed actions when clicked
			//
			// TODO: auto centering of action content.
			// Atm insets are hard-coded, thus actions have to know the bar height
			// and if bar height changes then actions would need to change.
			for _, action := range actions {
				action := action
				items = append(items, layout.Flexed(1, action))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, items...)
		}),
	)
}
