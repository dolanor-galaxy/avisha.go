package style

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

// Select implements the Material Design Select.
// described here: https://material-components.github.io/material-components-web-catalog/#/component/select
type Select struct {
	widget.Enum

	// Open if options are being rendered.
	Open bool

	values []string
	hint   string
	theme  *material.Theme

	options layout.List
	input   TextField
	toggle  widget.Clickable
	hovers  []widget.Hoverable
}

func (lb *Select) Update(gtx Ctx, th *material.Theme, hint string, values []string) {
	lb.theme = th
	lb.values = values
	lb.hint = hint
	if lb.options.Axis != layout.Vertical {
		lb.options.Axis = layout.Vertical
	}
	if lb.toggle.Clicked() {
		lb.Open = !lb.Open
	}
	if !lb.input.Focused() {
		lb.Open = false
	}
	if lb.Enum.Changed() {
		fmt.Printf("selected: %v\n", lb.Enum.Value)
		lb.input.Editor.SetText(lb.Enum.Value)
		lb.Open = false
	}
}

func (lb *Select) Layout(gtx Ctx, th *material.Theme, hint string, values []string) util.Dimensions {
	dims := lb.layout(gtx, th, hint, values)
	macro := op.Record(gtx.Ops)
	lb.Overlay(gtx)
	overlay := macro.Stop()
	return util.Dimensions{
		Dimensions: dims,
		Overlay:    overlay,
	}
}

func (lb *Select) layout(gtx Ctx, th *material.Theme, hint string, values []string) Dims {
	lb.Update(gtx, th, hint, values)
	return layout.Stack{}.Layout(
		gtx,
		layout.Stacked(func(gtx Ctx) Dims {
			return lb.input.Layout(gtx, th, hint)
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(
				gtx,
				layout.Flexed(1, func(gtx Ctx) Dims {
					return Dims{Size: image.Point{
						X: gtx.Constraints.Max.X,
						Y: gtx.Constraints.Min.Y,
					}}
				}),
				layout.Rigid(func(gtx Ctx) Dims {
					size := th.TextSize.Scale(2)
					return layout.Inset{
						Top:   unit.Px(float32(gtx.Px(size))/2 - float32(gtx.Px(unit.Dp(2)))),
						Right: unit.Dp(10),
					}.Layout(gtx, func(gtx Ctx) Dims {
						// TODO: animate icon with rotation.
						var arrow *widget.Icon
						if lb.Open {
							arrow = icons.ArrowUp
						} else {
							arrow = icons.ArrowDown
						}
						arrow.Color = lb.input.Color
						return arrow.Layout(gtx, size)
					})
				}),
			)
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			defer op.Push(gtx.Ops).Pop()
			pointer.PassOp{Pass: true}.Add(gtx.Ops)
			return lb.toggle.Layout(gtx)
		}),
	)
}

func (lb *Select) Overlay(gtx Ctx) Dims {
	var (
		values = lb.values
		th     = lb.theme
	)
	if !lb.Open {
		return Dims{}
	}
	return layout.Stack{}.Layout(
		gtx,
		// Debug color.
		layout.Expanded(func(gtx Ctx) Dims {
			return util.Rect{
				Color: color.RGBA{R: 0, G: 255, B: 255, A: 255},
				// Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
				Size: image.Point{
					X: gtx.Constraints.Max.X,
					Y: gtx.Constraints.Min.Y,
				},
				Radii: 4,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx Ctx) Dims {
			if len(lb.hovers) < len(values) {
				lb.hovers = make([]widget.Hoverable, len(values))
			}
			return layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(5),
			}.Layout(
				gtx,
				func(gtx Ctx) Dims {
					return lb.options.Layout(gtx, len(values), func(gtx Ctx, ii int) Dims {
						var (
							value = values[ii]
							hover = &lb.hovers[ii]
						)
						return layout.Stack{}.Layout(
							gtx,
							layout.Expanded(func(gtx Ctx) Dims {
								if hover.Hovered() {
									return util.Rect{
										Color: materials.AlphaMultiply(th.Color.Hint, 38),
										Size: image.Point{
											X: gtx.Constraints.Max.X,
											Y: gtx.Constraints.Min.Y,
										},
									}.Layout(gtx)
								}
								return Dims{
									Size: image.Point{
										X: gtx.Constraints.Max.X,
										Y: gtx.Constraints.Min.Y,
									},
								}
							}),
							layout.Stacked(func(gtx Ctx) Dims {
								return layout.UniformInset(unit.Dp(10)).Layout(
									gtx,
									func(gtx Ctx) Dims {
										dims := material.Label(th, th.TextSize, value).Layout(gtx)
										dims.Size.X = gtx.Constraints.Max.X
										return dims
									},
								)
							}),
							layout.Expanded(func(gtx Ctx) Dims {
								return lb.Enum.Layout(gtx, value)
							}),
							layout.Expanded(func(gtx Ctx) Dims {
								return hover.Layout(gtx)
							}),
						)
					})
				},
			)
		}),
	)
}
