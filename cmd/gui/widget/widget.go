// Package widget implements state for visual components.
// Re-exports some types from "gioui.org/widget".
package widget

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/util"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Re-export widget types.
type (
	Editor      = widget.Editor
	EditorEvent = widget.EditorEvent
	Enum        = widget.Enum
	Clickable   = widget.Clickable
	Bool        = widget.Bool
	Border      = widget.Border
	Icon        = widget.Icon
)

// Div is a visual divider: a colored line with a thickness.
type Div struct {
	Thickness unit.Value
	Length    unit.Value
	Axis      layout.Axis
	Color     color.NRGBA
}

func (d Div) Layout(gtx C) D {
	// Draw a line as a very thin rectangle.
	var sz image.Point
	switch d.Axis {
	case layout.Horizontal:
		sz = image.Point{
			X: gtx.Px(d.Length),
			Y: gtx.Px(d.Thickness),
		}
	case layout.Vertical:
		sz = image.Point{
			X: gtx.Px(d.Thickness),
			Y: gtx.Px(d.Length),
		}
	}
	return util.DrawRect(gtx, d.Color, sz, unit.Dp(0))
}
