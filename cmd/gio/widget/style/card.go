package style

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

// Card is a coherent stack of information with different content segments.
// Often used with header-body-footer style composition.
type Card struct {
	Content []layout.Widget
}

func (c Card) Layout(gtx C, th *material.Theme) D {
	var items = make([]layout.FlexChild, len(c.Content))
	for ii := range c.Content {
		ii := ii
		items[ii] = layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(5)).Layout(
				gtx,
				func(gtx C) D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							if ii == 0 {
								return D{}
							}
							return layout.Inset{
								Bottom: unit.Dp(10),
							}.Layout(
								gtx,
								func(gtx C) D {
									return widget.Div{
										Thickness: unit.Dp(1),
										Length:    unit.Px(float32(gtx.Constraints.Max.X)),
										Axis:      layout.Horizontal,
										Color:     color.RGBA{A: 100},
									}.Layout(gtx)
								},
							)
						}),
						layout.Rigid(func(gtx C) D {
							return c.Content[ii](gtx)
						}),
					)
				},
			)
		})
	}
	return widget.Border{
		Color:        th.Color.Hint,
		CornerRadius: unit.Dp(4),
		Width:        unit.Dp(0.5),
	}.Layout(
		gtx,
		func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx C) D {
					return util.DrawRect(
						gtx,
						color.RGBA{R: 255, G: 255, B: 255, A: 255},
						gtx.Constraints.Min,
						unit.Dp(4),
					)
				}),
				layout.Stacked(func(gtx C) D {
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
				}),
			)
		},
	)
}
