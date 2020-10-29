package theme

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/widget/material"
)

// Theme contains a semantic color palette.
type Theme struct {
	*material.Theme
	Palette
}

// Palette contains all semantic colors of a theme.
type Palette struct {
	Primary   color.RGBA
	Secondary color.RGBA
	Success   color.RGBA
	Info      color.RGBA
	Warning   color.RGBA
	Danger    color.RGBA
	Light     color.RGBA
	Dark      color.RGBA
}

// Option can be used to initialise a theme.
type Option func(*Theme)

// New theme instance.
func New(options ...Option) *Theme {
	th := &Theme{
		Theme: material.NewTheme(gofont.Collection()),
	}
	for _, opt := range options {
		opt(th)
	}
	return th
}

func (th Theme) Primary() *material.Theme {
	return with(th.Theme, th.Palette.Primary)
}

func (th Theme) Secondary() *material.Theme {
	return with(th.Theme, th.Palette.Secondary)
}

func (th Theme) Success() *material.Theme {
	return with(th.Theme, th.Palette.Success)
}

func (th Theme) Info() *material.Theme {
	return with(th.Theme, th.Palette.Info)
}

func (th Theme) Warning() *material.Theme {
	return with(th.Theme, th.Palette.Warning)
}

func (th Theme) Danger() *material.Theme {
	return with(th.Theme, th.Palette.Danger)
}

func (th Theme) Light() *material.Theme {
	return with(th.Theme, th.Palette.Light)
}

func (th Theme) Dark() *material.Theme {
	return with(th.Theme, th.Palette.Dark)
}

func with(base *material.Theme, c color.RGBA) *material.Theme {
	if base == nil {
		base = material.NewTheme(gofont.Collection())
	}
	th := *base
	th.Color.Primary = c
	return &th
}

// BootstrapPallet initialises theme with standard colors from bootstrap:
// https://getbootstrap.com/docs/4.0/utilities/colors/
func BootstrapPalette(th *Theme) {
	th.Palette = Palette{
		Primary:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		Secondary: color.RGBA{R: 108, G: 117, B: 125, A: 255},
		Success:   color.RGBA{R: 40, G: 167, B: 69, A: 255},
		Warning:   color.RGBA{R: 255, G: 193, B: 7, A: 255},
		Danger:    color.RGBA{R: 220, G: 53, B: 69, A: 255},
		Info:      color.RGBA{R: 23, G: 162, B: 184, A: 255},
		Light:     color.RGBA{R: 248, G: 249, B: 250, A: 255},
		Dark:      color.RGBA{R: 52, G: 58, B: 64, A: 255},
	}
}
