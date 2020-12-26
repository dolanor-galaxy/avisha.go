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

type TenantForm struct {
	nav.Route
	App *avisha.App
	Th  *style.Theme

	Tenant avisha.Tenant

	Name    materials.TextField
	Contact materials.TextField

	// @Improvemnt mixing form implementations is janky.
	Address AddressForm

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (f *TenantForm) Title() string {
	return "Tenant Form"
}

func (f *TenantForm) Receive(data interface{}) {
	if tenant, ok := data.(*avisha.Tenant); ok && tenant != nil {
		f.Tenant = *tenant
	} else {
		f.Tenant = avisha.Tenant{}
	}
	if f.Tenant.Address == (avisha.Address{}) {
		settings, err := f.App.LoadSettings()
		if err != nil {
			log.Printf("loading settings: %v", err)
		}
		f.Tenant.Address = settings.Defaults.Address
	}
	f.Address.Load(&f.Tenant.Address)
	f.Form.Load([]widget.Field{
		{
			Value: widget.RequiredValuer{Valuer: widget.TextValuer{Value: &f.Tenant.Name}},
			Input: &f.Name,
		},
		{
			Value: widget.TextValuer{Value: &f.Tenant.Contact},
			Input: &f.Contact,
		},
	})
}

func (f *TenantForm) Context() (list []layout.Widget) {
	if f.Tenant != (avisha.Tenant{}) {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(f.Th.Dark(), unit.Dp(24), f.Tenant.Name)
					label.Alignment = text.Middle
					label.Color = f.Th.Dark().ContrastFg
					return label.Layout(gtx)
				})
		})
	}
	return list
}

// Submit validates the input adata and returns a boolean indicating validity.
func (f *TenantForm) Submit() (tenant avisha.Tenant, ok bool) {
	return f.Tenant, f.Form.Submit()
}

func (f *TenantForm) Update(gtx C) {
	f.Form.Validate(gtx)
	if f.SubmitBtn.Clicked() {
		if t, ok := f.Submit(); ok {
			if err := func() error {
				if create := t.ID == 0; create {
					if err := f.App.RegisterTenant(&t); err != nil {
						return fmt.Errorf("registering tenant: %w", err)
					}
				} else {
					if err := f.App.Update(&t); err != nil {
						return fmt.Errorf("updating tenant: %w", err)
					}
					// Allow for zero value contact field.
					if err := f.App.UpdateField(
						&avisha.Tenant{ID: f.Tenant.ID},
						"Contact",
						t.Contact,
					); err != nil {
						return fmt.Errorf("updating tenant: %w", err)
					}
				}
				return nil
			}(); err != nil {
				log.Printf("%v", err)
			} else {
				f.Form.Clear()
				f.Route.Back()
			}
		}
	}
	if f.CancelBtn.Clicked() {
		f.Form.Clear()
		f.Route.Back()
	}
}

func (f *TenantForm) Layout(gtx C) D {
	f.Update(gtx)
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
						return f.Name.Layout(gtx, f.Th.Dark(), "Name")
					}),
					layout.Rigid(func(gtx C) D {
						return f.Contact.Layout(gtx, f.Th.Dark(), "Contact")
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return f.Address.Layout(gtx, f.Th)
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
								return material.Button(f.Th.Secondary(), &f.CancelBtn, "Cancel").Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
							}),
							layout.Rigid(func(gtx C) D {
								return material.Button(f.Th.Primary(), &f.SubmitBtn, "Submit").Layout(gtx)
							}),
						)
					},
				)
			}),
		)
	})
}
