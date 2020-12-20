package views

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
)

// LeaseForm performs data mutations on a Lease entity.
type LeaseForm struct {
	App *avisha.App

	Lease avisha.Lease

	// Form fields.
	Tenant materials.TextField
	Site   materials.TextField
	Date   materials.TextField
	Days   materials.TextField
	Rent   materials.TextField

	// Actions.
	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

// Submit validates the input data and returns a boolean indicating validity.
func (l *LeaseForm) Submit() (lease avisha.Lease, ok bool) {
	return l.Lease, l.Form.Submit()
}

func (l *LeaseForm) Clear() {
	l.Form.Clear()
}

// Load form data from a lease entity.
func (l *LeaseForm) Load(lease avisha.Lease) {
	l.Lease = lease
	l.Form.Load([]widget.Field{
		{
			Value: TenantValuer{ID: &l.Lease.Tenant, App: l.App},
			Input: &l.Tenant,
		},
		{
			Value: SiteValuer{ID: &l.Lease.Site, App: l.App},
			Input: &l.Site,
		},
		{
			Value: widget.DateValuer{Value: &l.Lease.Term.Start},
			Input: &l.Date,
		},
		{
			Value: widget.DaysValuer{Value: &l.Lease.Term.Duration},
			Input: &l.Days,
		},
		{
			Value: widget.CurrencyValuer{Value: &l.Lease.Rent},
			Input: &l.Rent,
		},
	})
}

func (l *LeaseForm) Layout(gtx C, th *style.Theme) D {
	l.Form.Validate(gtx)
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
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Flexed(1, func(gtx C) D {
							return l.Tenant.Layout(gtx, th.Dark(), "Tenant")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							return l.Site.Layout(gtx, th.Dark(), "Site")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return l.Date.Layout(gtx, th.Dark(), "Start Date")
				}),
				layout.Rigid(func(gtx C) D {
					l.Days.Suffix = func(gtx C) D {
						return material.Body1(th.Muted(), " days").Layout(gtx)
					}
					return l.Days.Layout(gtx, th.Dark(), "Duration")
				}),
				layout.Rigid(func(gtx C) D {
					l.Rent.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return l.Rent.Layout(gtx, th.Dark(), "Rent (weekly)")
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
							return material.Button(th.Secondary(), &l.CancelBtn, "Cancel").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							// if l.lease == nil {
							// 	return material.Button(th.Success(), &l.SubmitBtn, "Create").Layout(gtx)
							// }
							return material.Button(th.Primary(), &l.SubmitBtn, "Update").Layout(gtx)
						}),
					)
				})
		}),
	)
}
