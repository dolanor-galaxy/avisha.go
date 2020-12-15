package style

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
)

// Card is a coherent stack of information with different content segments.
// Often used with header-body-footer style composition.
type Card struct {
	Content []layout.Widget
}

func (c Card) Layout(gtx C, th *material.Theme) D {
	var (
		items = make([]layout.FlexChild, len(c.Content))
	)
	for ii := range c.Content {
		ii := ii
		items[ii] = layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(5)).Layout(
				gtx,
				func(gtx C) D {
					return c.Content[ii](gtx)
				},
			)
		})
	}
	return layout.Stack{}.Layout(
		gtx,
		layout.Expanded(func(gtx C) D {
			return util.DrawRect(
				gtx,
				color.NRGBA{R: 255, G: 255, B: 255, A: 255},
				gtx.Constraints.Min,
				unit.Dp(4),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return widget.Border{
				Color:        th.ContrastBg,
				CornerRadius: unit.Dp(4),
				Width:        unit.Dp(0.5),
			}.Layout(
				gtx,
				func(gtx C) D {
					return layout.UniformInset(unit.Dp(5)).Layout(
						gtx,
						func(gtx C) D {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(
								gtx,
								items...,
							)
						},
					)
				},
			)
		}),
	)
}
