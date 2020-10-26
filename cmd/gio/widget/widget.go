// Package widget implements state for visual components.
// Re-exports some types from "gioui.org/widget".
package widget

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Re-export widget types.
type (
	Editor    = widget.Editor
	Enum      = widget.Enum
	Clickable = widget.Clickable
	Bool      = widget.Bool
	Border    = widget.Border
	Icon      = widget.Icon
)
