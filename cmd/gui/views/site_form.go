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
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
)

type SiteForm struct {
	// Page state.
	nav.Route
	App *avisha.App
	Th  *style.Theme

	// Entity data.
	site *avisha.Site

	// Form field.
	Number materials.TextField

	// Actions.
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
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
					label := material.Label(l.Th.Primary(), unit.Dp(24), l.site.Number)
					label.Alignment = text.Middle
					label.Color = l.Th.Primary().Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

// Clear the form fields.
func (l *SiteForm) Clear() {
	l.Number.Clear()
	l.site = nil
}

func (l *SiteForm) Update(gtx C) {
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
				l.Clear()
				l.Route.Back()
			}
		}
	}
	if l.CancelBtn.Clicked() {
		l.Clear()
		l.Route.Back()
	}
}

func (l *SiteForm) Layout(gtx C) D {
	l.Update(gtx)
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
		return style.Container{
			BreakPoint: unit.Dp(700),
			Constrain:  true,
		}.Layout(gtx, func(gtx C) D {
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
	})
}

// Submit validates form data and returns a boolean to indicate validity.
func (l *SiteForm) Submit() (s avisha.Site, ok bool) {
	ok = true
	if l.site != nil {
		s.ID = l.site.ID
	}
	if n, err := util.FieldRequired(l.Number.Text()); err != nil {
		l.Number.SetError(err.Error())
		ok = false
	} else {
		s.Number = n
	}
	return s, ok
}
