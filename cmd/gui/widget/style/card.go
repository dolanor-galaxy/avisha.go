package style

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
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
		items   = make([]layout.FlexChild, len(c.Content))
		calls   = make([]op.CallOp, len(c.Content))
		dimlist = make([]D, len(c.Content))
		width   = 0
	)
	// First pass calculates max width of the card so we can size the divs later.
	for ii := range c.Content {
		macro := op.Record(gtx.Ops)
		dims := c.Content[ii](gtx)
		call := macro.Stop()
		if dims.Size.X > width {
			width = dims.Size.X
		}
		calls[ii] = call
		dimlist[ii] = dims
	}
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
							if skipFirst := ii == 0; skipFirst {
								return D{}
							}
							return layout.Inset{
								Bottom: unit.Dp(10),
							}.Layout(
								gtx,
								func(gtx C) D {
									return widget.Div{
										Thickness: unit.Dp(1),
										Length:    unit.Px(float32(width)),
										Axis:      layout.Horizontal,
										Color:     color.NRGBA{A: 100},
									}.Layout(gtx)
								},
							)
						}),
						layout.Rigid(func(gtx C) D {
							// return c.Content[ii](gtx)
							calls[ii].Add(gtx.Ops)
							return dimlist[ii]
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
						color.NRGBA{R: 255, G: 255, B: 255, A: 255},
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
