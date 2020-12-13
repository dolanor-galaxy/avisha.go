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

type TenantForm struct {
	// Page data.
	nav.Route
	App *avisha.App
	Th  *style.Theme

	// Entity data.
	tenant *avisha.Tenant

	// Form fields.
	Name    materials.TextField
	Contact materials.TextField

	// Actions.
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (f *TenantForm) Title() string {
	return "Tenant Form"
}

func (f *TenantForm) Receive(data interface{}) {
	if tenant, ok := data.(*avisha.Tenant); ok && tenant != nil {
		f.tenant = tenant
		f.Name.SetText(tenant.Name)
		f.Contact.SetText(tenant.Contact)
	}
}

func (f *TenantForm) Context() (list []layout.Widget) {
	if f.tenant != nil {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(f.Th.Primary(), unit.Dp(24), f.tenant.Name)
					label.Alignment = text.Middle
					label.Color = f.Th.Primary().Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

// Clear the form fields.
func (f *TenantForm) Clear() {
	f.Name.Clear()
	f.Contact.Clear()
	f.tenant = nil
}

func (f *TenantForm) Update(gtx C) {
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
						&avisha.Tenant{ID: f.tenant.ID},
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
				f.Clear()
				f.Route.Back()
			}
		}
	}
	if f.CancelBtn.Clicked() {
		f.Clear()
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
						return f.Name.Layout(gtx, f.Th.Primary(), "Name")
					}),
					layout.Rigid(func(gtx C) D {
						return f.Contact.Layout(gtx, f.Th.Primary(), "Contact")
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

// Submit validates the input adata and returns a boolean indicating validity.
func (f *TenantForm) Submit() (tenant avisha.Tenant, ok bool) {
	ok = true
	if f.tenant != nil {
		tenant.ID = f.tenant.ID
	}
	if name, err := f.validateName(); err != nil {
		f.Name.SetError(err.Error())
		ok = false
	} else {
		tenant.Name = name
	}
	tenant.Contact = f.Contact.Text()
	return tenant, ok
}

func (f *TenantForm) validateName() (string, error) {
	return util.FieldRequired(f.Name.Text())
}
