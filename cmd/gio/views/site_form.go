package views

import (
	"fmt"
	"image"
	"log"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/theme"
)

type SiteForm struct {
	nav.Route
	App  *avisha.App
	Th   *theme.Theme
	site *avisha.Site

	Number materials.TextField
	Submit widget.Clickable
	Cancel widget.Clickable
}

func (l *SiteForm) Title() string {
	return "Site Form"
}

func (l *SiteForm) Receive(data interface{}) {
	if site, ok := data.(*avisha.Site); ok && site != nil {
		l.site = site
		l.Number.SetText(site.Number)
	}
}

func (l *SiteForm) Context() (list []layout.Widget) {
	if l.site != nil {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(l.Th.Primary(), unit.Dp(24), l.site.ID())
					label.Alignment = text.Middle
					label.Color = l.Th.Primary().Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (l *SiteForm) Update(gtx C) {
	if l.Submit.Clicked() {
		if err := l.submit(); err != nil {
			// give error to app or render under field.
			log.Printf("listing site form: %v", err)
		}
		l.Route.Back()
	}
	if l.Cancel.Clicked() {
		l.Route.Back()
	}
}

func (l *SiteForm) Layout(gtx C) D {
	l.Update(gtx)
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx C) D {
					l.Number.SingleLine = true
					return l.Number.Layout(gtx, l.Th.Primary(), "Number")
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top: unit.Dp(10),
			}.Layout(
				gtx,
				func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							return material.Button(l.Th.Secondary(), &l.Cancel, "Cancel").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							return material.Button(l.Th.Primary(), &l.Submit, "Submit").Layout(gtx)
						}),
					)
				})
		}),
	)
}

func (l *SiteForm) submit() error {
	site := avisha.Site{
		Number: l.Number.Text(),
		// Dwelling: l.Dwelling.Text(),
		Dwelling: avisha.Cabin,
	}
	if err := l.App.ListSite(site); err != nil {
		return fmt.Errorf("listing site: %w", err)
	}
	return nil
}
