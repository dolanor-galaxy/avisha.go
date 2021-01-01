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

// SettingsForm performs manipulations of settings.
type SettingsForm struct {
	Settings *avisha.Settings

	// @Todo(low) user-driven fields like BSB for Australian banks.
	Bank struct {
		Name    materials.TextField
		Account materials.TextField
	}

	Landlord struct {
		Name  materials.TextField
		Email materials.TextField
		Phone materials.TextField
		AddressForm
	}

	Defaults struct {
		UnitCost   materials.TextField
		RentCycle  materials.TextField
		InvoiceNet materials.TextField
		GST        materials.TextField
	}

	// // BillTo is the default billable address for Tenants.
	// BillTo struct {
	// 	Unit   materials.TextField
	// 	Number materials.TextField
	// 	Street materials.TextField
	// 	City   materials.TextField
	// }

	BillTo AddressForm

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (s *SettingsForm) Clear() {
	if s.Settings != nil {
		s.Load(s.Settings)
	} else {
		s.Form.Clear()
	}
}

// Load initialises the form fields.
func (s *SettingsForm) Load(settings *avisha.Settings) {
	s.Settings = settings
	s.Landlord.Load(&s.Settings.Landlord.Address)
	s.BillTo.Load(&s.Settings.Defaults.Address)
	s.Form.Load([]widget.Field{
		{
			Value: widget.TextValuer{Value: &s.Settings.Landlord.Name},
			Input: &s.Landlord.Name,
		},
		{
			Value: widget.TextValuer{Value: &s.Settings.Landlord.Email},
			Input: &s.Landlord.Email,
		},
		{
			Value: widget.TextValuer{Value: &s.Settings.Landlord.Phone},
			Input: &s.Landlord.Phone,
		},
		{
			Value: widget.TextValuer{Value: &s.Settings.Bank.Name},
			Input: &s.Bank.Name,
		},
		{
			Value: widget.TextValuer{Value: &s.Settings.Bank.Account},
			Input: &s.Bank.Account,
		},
		{
			Value: widget.CurrencyValuer{Value: &s.Settings.Defaults.UnitCost},
			Input: &s.Defaults.UnitCost,
		},
		{
			Value: widget.DaysValuer{Value: &s.Settings.Defaults.RentCycle},
			Input: &s.Defaults.RentCycle,
		},
		{
			Value: widget.DaysValuer{Value: &s.Settings.Defaults.InvoiceNet},
			Input: &s.Defaults.InvoiceNet,
		},
		{
			Value: widget.FloatValuer{Value: &s.Settings.Defaults.GST},
			Input: &s.Defaults.GST,
		},
	})
}

// Submit validates the data and returns a boolean indicating validity.
func (s *SettingsForm) Submit() (settings avisha.Settings, ok bool) {
	if !s.Landlord.Form.Submit() {
		return settings, false
	}
	if !s.BillTo.Form.Submit() {
		return settings, false
	}
	if !s.Form.Submit() {
		return settings, false
	}
	return *s.Settings, true
}

func (s *SettingsForm) Layout(gtx C, th *style.Theme) D {
	spacer := func(size ...unit.Value) layout.Widget {
		return func(gtx C) D {
			var sz int
			if len(size) > 0 {
				for ii := range size {
					sz += gtx.Px(size[ii])
				}
			} else {
				sz = gtx.Px(unit.Dp(10))
			}
			return D{Size: image.Point{X: sz, Y: sz}}
		}
	}
	title := func(title string, size ...unit.Value) layout.Widget {
		var (
			top    = unit.Dp(20)
			bottom = unit.Dp(10)
		)
		if len(size) > 0 {
			top = size[0]
		}
		if len(size) > 1 {
			bottom = size[1]
		}
		return func(gtx C) D {
			return layout.Inset{
				Top:    top,
				Bottom: bottom,
			}.Layout(gtx, func(gtx C) D {
				return material.Label(th.Dark(), unit.Dp(20), title).Layout(gtx)
			})
		}
	}
	field := func(f *materials.TextField, name string, options ...func(f *materials.TextField)) layout.Widget {
		return func(gtx C) D {
			for _, opt := range options {
				opt(f)
			}
			return f.Layout(gtx, th.Dark(), name)
		}
	}
	s.Form.Validate(gtx)
	s.BillTo.Form.Validate(gtx)
	s.Landlord.AddressForm.Form.Validate(gtx)
	return layout.UniformInset(unit.Dp(10)).Layout(
		gtx,
		func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,
				layout.Rigid(title("Landlord", unit.Dp(0))),
				layout.Rigid(field(&s.Landlord.Name, "Name")),
				layout.Rigid(field(&s.Landlord.Email, "Email")),
				layout.Rigid(field(&s.Landlord.Phone, "Phone")),
				layout.Rigid(title("Landlord Address", unit.Dp(10))),
				layout.Rigid(func(gtx C) D {
					return s.Landlord.AddressForm.Layout(gtx, th)
				}),
				layout.Rigid(title("Bank Details")),
				layout.Rigid(field(&s.Bank.Name, "Name")),
				layout.Rigid(field(&s.Bank.Account, "Account")),
				layout.Rigid(title("Defaults")),
				layout.Rigid(field(
					&s.Defaults.UnitCost,
					"Unit Cost (dollars)",
					func(f *materials.TextField) {
						f.Prefix = func(gtx C) D {
							return material.Body1(th.Theme, "$").Layout(gtx)
						}
					})),
				layout.Rigid(field(&s.Defaults.RentCycle, "Rent Cycle (days)")),
				layout.Rigid(field(&s.Defaults.InvoiceNet, "Invoice Net (days)")),
				layout.Rigid(field(
					&s.Defaults.GST,
					"GST",
					func(f *materials.TextField) {
						f.Suffix = func(gtx C) D {
							return material.Body1(th.Theme, "%").Layout(gtx)
						}
					})),
				layout.Rigid(title("Default Billable Address")),
				layout.Rigid(func(gtx C) D {
					return s.BillTo.Layout(gtx, th)
				}),
				layout.Rigid(spacer(unit.Dp(20))),
				// @Todo sticky controls.
				layout.Rigid(func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							return material.Button(th.Muted(), &s.CancelBtn, "Cancel").Layout(gtx)
						}),
						layout.Rigid(spacer()),
						layout.Rigid(func(gtx C) D {
							return material.Button(th.Primary(), &s.SubmitBtn, "Update").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return D{Size: gtx.Constraints.Min}
						}),
					)
				}),
			)
		},
	)
}
