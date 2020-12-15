package style

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
	Primary   material.Palette
	Secondary material.Palette
	Success   material.Palette
	Info      material.Palette
	Warning   material.Palette
	Danger    material.Palette
	Light     material.Palette
	Dark      material.Palette
	Muted     material.Palette
}

// ThemeOption can be used to initialise a theme.
type ThemeOption = func(*Theme)

// NewTheme allocates a theme instance.
func NewTheme(options ...ThemeOption) *Theme {
	th := &Theme{
		Theme: material.NewTheme(gofont.Collection()),
	}
	for _, opt := range options {
		opt(th)
	}
	return th
}

func (th Theme) Primary() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Primary)
	return &mth
}

func (th Theme) Secondary() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Secondary)
	return &mth
}

func (th Theme) Success() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Success)
	return &mth
}

func (th Theme) Info() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Info)
	return &mth
}

func (th Theme) Warning() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Warning)
	return &mth
}

func (th Theme) Danger() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Danger)
	return &mth
}

func (th Theme) Light() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Light)
	return &mth
}

func (th Theme) Dark() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Dark)
	return &mth
}

func (th Theme) Muted() *material.Theme {
	mth := th.Theme.WithPalette(th.Palette.Muted)
	return &mth
}

// BootstrapPalette initialises theme with standard colors from bootstrap:
// https://getbootstrap.com/docs/4.0/utilities/colors/
func BootstrapPalette(th *Theme) {
	th.Palette = Palette{
		Primary:   Light(rgb(0x007bff)),
		Secondary: Light(rgb(0x6c757d)),
		Success:   Light(rgb(0x28a745)),
		Danger:    Light(rgb(0xdc3545)),
		Warning:   Light(rgb(0xffc107)),
		Info:      Light(rgb(0x17a2b8)),
		Light:     Light(rgb(0xf8f9fa)),
		Dark:      Light(rgb(0x343a40)),
		Muted:     Light(rgb(0x6c757d)),
	}
}

func BootstrapDarkPalette(th *Theme) {
	th.Palette = Palette{
		Primary:   Dark(rgb(0x007bff)),
		Secondary: Dark(rgb(0x6c757d)),
		Success:   Dark(rgb(0x28a745)),
		Danger:    Dark(rgb(0xdc3545)),
		Warning:   Dark(rgb(0xffc107)),
		Info:      Dark(rgb(0x17a2b8)),
		Light:     Dark(rgb(0xf8f9fa)),
		Dark:      Dark(rgb(0x343a40)),
		Muted:     Dark(rgb(0x6c757d)),
	}
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func Light(c color.NRGBA) material.Palette {
	return material.Palette{
		Fg:         c,
		Bg:         rgb(0xffffff),
		ContrastBg: c,
		ContrastFg: rgb(0xffffff),
	}
}

func Dark(c color.NRGBA) material.Palette {
	return material.Palette{
		Fg:         c,
		Bg:         rgb(0x000000),
		ContrastBg: c,
		ContrastFg: rgb(0x000000),
	}
}
