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

type TenantForm struct {
	nav.Route
	App    *avisha.App
	Th     *theme.Theme
	tenant *avisha.Tenant

	Name    materials.TextField
	Contact materials.TextField
	Submit  widget.Clickable
	Cancel  widget.Clickable
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
					label := material.Label(f.Th.Primary(), unit.Dp(24), f.tenant.ID())
					label.Alignment = text.Middle
					label.Color = f.Th.Primary().Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (f *TenantForm) Update(gtx C) {
	if f.Submit.Clicked() {
		if err := f.submit(); err != nil {
			// give error to app or render under field.
			log.Printf("submitting tenant form: %v", err)
		}
		f.Route.Back()
	}
	if f.Cancel.Clicked() {
		f.Route.Back()
	}
}

func (f *TenantForm) Layout(gtx C) D {
	f.Update(gtx)
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
							return material.Button(f.Th.Secondary(), &f.Cancel, "Cancel").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							return material.Button(f.Th.Primary(), &f.Submit, "Submit").Layout(gtx)
						}),
					)
				},
			)
		}),
	)
}

func (f *TenantForm) submit() error {
	t := avisha.Tenant{
		Name:    f.Name.Text(),
		Contact: f.Contact.Text(),
	}
	if err := f.App.RegisterTenant(t); err != nil {
		return fmt.Errorf("registering tenant: %w", err)
	}
	return nil
}
