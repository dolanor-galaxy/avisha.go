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

// Container treats content differently based on viewport width.
type Container struct {
	// BreakPoint describes when to change the layout.
	BreakPoint unit.Value
	// Constrain limits widget width by the breakpoint.
	Constrain bool
	// Scroll allows verical scrolling.
	Scroll bool

	layout.List
}

// @Todo: implement scrolling. Using list is tricky because it simulates an
// infinite Y max.
func (c Container) Layout(gtx C, w layout.Widget) D {
	c.List.Axis = layout.Vertical
	breakpoint := gtx.Px(c.BreakPoint)
	// if c.Scroll {
	// 	return c.List.Layout(gtx, 1, func(gtx C, _ int) D {
	// 		return Centered(gtx, func(gtx C) D {
	// 			cs := &gtx.Constraints
	// 			if c.Constrain && cs.Max.X > breakpoint {
	// 				cs.Max.X = breakpoint
	// 			}
	// 			if c.Constrain && cs.Max.Y > breakpoint {
	// 				cs.Max.Y = breakpoint
	// 			}
	// 			return w(gtx)
	// 		})
	// 	})
	// }
	return Centered(gtx, func(gtx C) D {
		cs := &gtx.Constraints
		if c.Constrain && cs.Max.X > breakpoint {
			cs.Max.X = breakpoint
		}
		return w(gtx)
	})
}

func ModalDialog(gtx C, th *Theme, max unit.Value, title string, w layout.Widget) D {
	return Modal(gtx, max, func(gtx C) D {
		return Card{
			Content: []layout.Widget{
				func(gtx C) D {
					return material.Label(th.Primary(), unit.Dp(20), title).Layout(gtx)
				},
				func(gtx C) D {
					return w(gtx)
				},
			},
		}.Layout(gtx, th.Primary())
	})
}

// Modal renders content centered on a translucent scrim with a max width capped
// by the amount specified.
func Modal(gtx C, max unit.Value, w layout.Widget) D {
	return layout.Stack{}.Layout(
		gtx,
		layout.Stacked(func(gtx C) D {
			return util.Rect{
				Size:  gtx.Constraints.Max,
				Color: color.NRGBA{A: 200},
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			return Centered(gtx, func(gtx C) D {
				cs := &gtx.Constraints
				if cs.Max.X > gtx.Px(max) {
					cs.Max.X = gtx.Px(max)
				}
				return w(gtx)
			})
		}),
	)
}

// Centered places the widget in the center of the container.
func Centered(gtx C, w layout.Widget) D {
	return CenteredHorizontal(gtx, func(gtx C) D {
		return CenteredVertical(gtx, w)
	})
}

// CenteredHorizontal centers the widget along the horizontal axis.
func CenteredHorizontal(gtx C, w layout.Widget) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(
		gtx,
		layout.Flexed(1, func(gtx C) D {
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Rigid(func(gtx C) D {
			return w(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return D{Size: gtx.Constraints.Min}
		}),
	)
}

// CenteredVertical centers the widget along the vertical axis.
func CenteredVertical(gtx C, w layout.Widget) D {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(
		gtx,
		layout.Flexed(1, func(gtx C) D {
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Rigid(func(gtx C) D {
			return w(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return D{Size: gtx.Constraints.Min}
		}),
	)
}
