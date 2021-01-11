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

type UtilitiesPayment struct {
	InvoiceID int
	avisha.Payment
}

// UtilitiesPaymentForm allows data entry for paying utility invoices.
type UtilitiesPaymentForm struct {
	Payment UtilitiesPayment

	// @Note do we want to override payment date?
	// Probably
	Amount materials.TextField
	// @Todo use select
	Invoice materials.TextField

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (f *UtilitiesPaymentForm) Load(p UtilitiesPayment, invoices []avisha.Invoice) {
	f.Payment = p
	var (
		def = 0
	)
	if len(invoices) > 0 {
		def = invoices[len(invoices)-1].ID
	}
	f.Form.Load([]widget.Field{
		{
			Value: widget.IntValuer{
				Value: &f.Payment.InvoiceID,
				// default to oldest overdue invoice.
				Default: def,
			},
			Input: &f.Invoice,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Payment.Amount,
			},
			Input: &f.Amount,
		},
	})
}

func (f *UtilitiesPaymentForm) Submit() (p UtilitiesPayment, ok bool) {
	return f.Payment, f.Form.Submit()
}

func (f *UtilitiesPaymentForm) Clear() {
	f.Form.Clear()
}

func (f *UtilitiesPaymentForm) Layout(gtx C, th *style.Theme) D {
	f.Form.Validate(gtx)
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return f.Invoice.Layout(gtx, th.Muted(), "Invoice")
		}),
		layout.Rigid(func(gtx C) D {
			f.Amount.Prefix = func(gtx C) D {
				return material.Body2(th.Muted(), "$").Layout(gtx)
			}
			return f.Amount.Layout(gtx, th.Muted(), "Amount")
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
							return material.Button(th.Secondary(), &f.CancelBtn, "Cancel").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							return material.Button(th.Primary(), &f.SubmitBtn, "Create").Layout(gtx)
						}),
					)
				})
		}),
	)
}
