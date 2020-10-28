package nav

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
)

// Page wraps content in navigation context, such as a top bar and nav rail.
type Page struct {
	*material.Theme
	Router Router
	Rail   Rail
	layout.List
}

func (p *Page) Layout(gtx C) D {
	p.List.Axis = layout.Vertical
	for _, d := range p.Rail.Destinations {
		if d.Clicked() {
			p.Router.Push(d.Route, nil)
		}
	}
	for ii := range p.Rail.Destinations {
		d := &p.Rail.Destinations[ii]
		d.Active = d.Route == p.Router.Name()
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return style.TopBar{
				Theme:  p.Theme,
				Height: unit.Dp(50),
			}.Layout(
				gtx,
				func() string {
					if titled, ok := p.Router.Active().(Titled); ok {
						return titled.Title()
					}
					return ""
				}(),
				func() []layout.Widget {
					if contexter, ok := p.Router.Active().(Contexter); ok {
						return contexter.Context()
					}
					return nil
				}()...)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx C) D {
					return p.Rail.Layout(gtx, p.Theme)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.UniformInset(unit.Dp(10)).Layout(
						gtx,
						func(gtx C) D {
							return p.Router.Layout(gtx)
						},
					)
					// FIXME: nested lists do not scroll: how to scroll both list and page?
					// return p.List.Layout(gtx, 1, func(gtx C, _ int) D {
					// 	return p.Router.Layout(gtx)
					// })
				}),
			)
		}),
	)
}
