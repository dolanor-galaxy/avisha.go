package views

import (
	"image"
	"log"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/nav"
	"github.com/jackmordaunt/avisha.go/cmd/gui/util"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
	"github.com/jackmordaunt/avisha.go/currency"
)

// SettingsPage is a page for configuring global settings.
type SettingsPage struct {
	nav.Route
	App  *avisha.App
	Th   *style.Theme
	Form SettingsForm

	scroll layout.List
}

func (s *SettingsPage) Title() string {
	return "Settings"
}

func (s *SettingsPage) Receive(_ interface{}) {
	s.Load()
}

func (s *SettingsPage) Load() {
	if settings, err := s.App.LoadSettings(); err != nil {
		log.Printf("loading settings: %v", err)
	} else {
		s.Form.Load(&settings)
	}
}

func (s *SettingsPage) Update(gtx C) {
	if s.Form.SubmitBtn.Clicked() {
		if settings, ok := s.Form.Submit(); ok {
			if err := s.App.SaveSettings(settings); err != nil {
				log.Printf("updating settings: %v", err)
			}
		}
	}
	if s.Form.CancelBtn.Clicked() {
		s.Load()
	}
}

func (s *SettingsPage) Layout(gtx C) D {
	s.scroll.Axis = layout.Vertical
	s.Update(gtx)
	return s.scroll.Layout(gtx, 1, func(gtx C, index int) D {
		return s.Form.Layout(gtx, s.Th)
	})
}

// SettingsForm performs manipulations of settings.
type SettingsForm struct {
	Settings *avisha.Settings

	Bank
	Landlord
	Defaults

	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

// Landlord details.
type Landlord struct {
	Name  materials.TextField
	Email materials.TextField
	Phone materials.TextField
}

// Bank details to make invoices payable to.
// @Todo(low) user-driven arbitrary fields like BSB for Australian banks.
// Generalise an abstraction, such as dynamic form?
type Bank struct {
	Name    materials.TextField
	Account materials.TextField
}

// Defaults spcecifies global defaults to auto populate fields with.
type Defaults struct {
	UnitCost   materials.TextField
	RentCycle  materials.TextField
	InvoiceNet materials.TextField
}

func (s *SettingsForm) Layout(gtx C, th *style.Theme) D {
	s.Update(gtx)
	spacer := func(size ...unit.Value) layout.Widget {
		return func(gtx C) D {
			var sz unit.Value
			if len(size) > 0 {
				for ii := range size {
					sz.V += size[ii].V
				}
			} else {
				sz = unit.Dp(10)
			}
			return D{Size: image.Point{X: gtx.Px(sz), Y: gtx.Px(sz)}}
		}
	}
	title := func(title string) layout.Widget {
		return func(gtx C) D {
			return material.Label(th.Dark(), unit.Dp(20), title).Layout(gtx)
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
	return layout.UniformInset(unit.Dp(10)).Layout(
		gtx,
		func(gtx C) D {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,
				layout.Rigid(title("Landlord")),
				layout.Rigid(field(&s.Landlord.Name, "Name")),
				layout.Rigid(field(&s.Landlord.Email, "Email")),
				layout.Rigid(field(&s.Landlord.Phone, "Phone")),
				layout.Rigid(spacer()),
				layout.Rigid(title("Bank Details")),
				layout.Rigid(field(&s.Bank.Name, "Name")),
				layout.Rigid(field(&s.Bank.Account, "Account")),
				layout.Rigid(spacer()),
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
				layout.Rigid(spacer()),
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

// @Improvement reduce this form field boilerplate and apply to other forms.
// - realtime validation
// - submission validation
// - load/clear, mapping between fields and concrete values
//
// Bi-directional mapping = text <-> primitive <-> structure
// 1. Declarative (declare this relationship once)
// 2. Handle realtime validation and submission validation usecases
func (s *SettingsForm) Update(gtx C) {
	for range s.Defaults.UnitCost.Events() {
		if _, err := s.validateUnitCost(); err != nil {
			s.Defaults.UnitCost.SetError(err.Error())
		} else {
			s.Defaults.UnitCost.ClearError()
		}
	}
	for range s.Defaults.InvoiceNet.Events() {
		if _, err := s.validateInvoiceNet(); err != nil {
			s.Defaults.InvoiceNet.SetError(err.Error())
		} else {
			s.Defaults.InvoiceNet.ClearError()
		}
	}
	for range s.Defaults.RentCycle.Events() {
		if _, err := s.validateRentCycle(); err != nil {
			s.Defaults.RentCycle.SetError(err.Error())
		} else {
			s.Defaults.RentCycle.ClearError()
		}
	}
}

func (s *SettingsForm) Clear() {
	if s.Settings == nil {
		s.Bank.Account.Clear()
		s.Bank.Name.Clear()
		s.Landlord.Name.Clear()
		s.Landlord.Email.Clear()
		s.Landlord.Phone.Clear()
		s.Defaults.InvoiceNet.Clear()
		s.Defaults.RentCycle.Clear()
		s.Defaults.UnitCost.Clear()
	} else {
		s.Load(s.Settings)
	}
}

func (s *SettingsForm) Load(settings *avisha.Settings) {
	s.Settings = settings
	s.Bank.Account.SetText(settings.Bank.Account)
	s.Bank.Name.SetText(settings.Bank.Name)
	s.Landlord.Name.SetText(settings.Landlord.Name)
	s.Landlord.Email.SetText(settings.Landlord.Email)
	s.Landlord.Phone.SetText(settings.Landlord.Phone)
	s.Defaults.InvoiceNet.SetText(strconv.Itoa(int(settings.Defaults.InvoiceNet.Hours() / 24)))
	s.Defaults.RentCycle.SetText(strconv.Itoa(int(settings.Defaults.RentCycle.Hours() / 24)))
	s.Defaults.UnitCost.SetText(strconv.Itoa(int(settings.Defaults.UnitCost)))
}

// Submit validates the data and returns a boolean indicating validity.
func (s *SettingsForm) Submit() (settings avisha.Settings, ok bool) {
	ok = true
	settings.Bank.Name = s.Bank.Name.Text()
	settings.Bank.Account = s.Bank.Account.Text()
	settings.Landlord.Name = s.Landlord.Name.Text()
	settings.Landlord.Email = s.Landlord.Email.Text()
	settings.Landlord.Phone = s.Landlord.Phone.Text()
	if n, err := s.validateUnitCost(); err != nil {
		s.Defaults.UnitCost.SetError(err.Error())
		ok = false
	} else {
		settings.Defaults.UnitCost = n
	}
	if d, err := s.validateRentCycle(); err != nil {
		s.Defaults.RentCycle.SetError(err.Error())
		ok = false
	} else {
		settings.Defaults.RentCycle = d
	}
	if d, err := s.validateInvoiceNet(); err != nil {
		s.Defaults.InvoiceNet.SetError(err.Error())
		ok = false
	} else {
		settings.Defaults.InvoiceNet = d
	}
	return settings, ok
}

func (s *SettingsForm) validateUnitCost() (currency.Currency, error) {
	return util.ParseCurrency(s.Defaults.UnitCost.Text())
}

func (s *SettingsForm) validateInvoiceNet() (time.Duration, error) {
	return util.ParseDay(s.Defaults.InvoiceNet.Text())
}

func (s *SettingsForm) validateRentCycle() (time.Duration, error) {
	return util.ParseDay(s.Defaults.RentCycle.Text())
}
