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
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/nav"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
)

type SiteForm struct {
	nav.Route
	App *avisha.App
	Th  *style.Theme

	Site avisha.Site

	Number materials.TextField

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (l *SiteForm) Title() string {
	return "Site Form"
}

func (l *SiteForm) Receive(data interface{}) {
	if site, ok := data.(*avisha.Site); ok && site != nil {
		l.Site = *site
	} else {
		l.Site = avisha.Site{}
	}
	l.Form.Load([]widget.Field{
		{
			Value: widget.RequiredValuer{Valuer: widget.TextValuer{Value: &l.Site.Number}},
			Input: &l.Number,
		},
	})
}

func (l *SiteForm) Context() (list []layout.Widget) {
	if l.Site != (avisha.Site{}) {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(l.Th.Dark(), unit.Dp(24), l.Site.Number)
					label.Alignment = text.Middle
					label.Color = l.Th.Dark().ContrastFg
					return label.Layout(gtx)
				})
		})
	}
	return list
}

// Submit validates form data and returns a boolean to indicate validity.
func (l *SiteForm) Submit() (s avisha.Site, ok bool) {
	return l.Site, l.Form.Submit()
}

func (l *SiteForm) Update(gtx C) {
	l.Form.Validate(gtx)
	if l.SubmitBtn.Clicked() {
		if s, ok := l.Submit(); ok {
			if err := func() error {
				if create := s.ID == 0; create {
					if err := l.App.ListSite(&s); err != nil {
						return fmt.Errorf("listing site: %w", err)
					}
				} else {
					if err := l.App.Update(&s); err != nil {
						return fmt.Errorf("updating site: %w", err)
					}
				}
				return nil
			}(); err != nil {
				log.Printf("%v", err)
			} else {
				l.Form.Clear()
				l.Route.Back()
			}
		}
	}
	if l.CancelBtn.Clicked() {
		l.Form.Clear()
		l.Route.Back()
	}
}

func (l *SiteForm) Layout(gtx C) D {
	l.Update(gtx)
	if breakpoint := gtx.Px(unit.Dp(700)); gtx.Constraints.Max.X > breakpoint {
		gtx.Constraints.Max.X = breakpoint
	}
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
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
						return l.Number.Layout(gtx, l.Th.Dark(), "Number")
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
								return material.Button(l.Th.Secondary(), &l.CancelBtn, "Cancel").Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
							}),
							layout.Rigid(func(gtx C) D {
								return material.Button(l.Th.Primary(), &l.SubmitBtn, "Submit").Layout(gtx)
							}),
						)
					})
			}),
		)
	})
}
