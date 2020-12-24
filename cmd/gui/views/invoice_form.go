package views

import (
	"image"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/util"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
	"github.com/jackmordaunt/avisha.go/currency"
)

// UtilitiesInvoiceForm is form for collecting utility invoice data.
type UtilitiesInvoiceForm struct {
	Invoice avisha.UtilityInvoice

	UnitsConsumed   materials.TextField
	PreviousReading materials.TextField
	CurrentReading  materials.TextField
	UnitCost        materials.TextField
	TotalCost       materials.TextField
	IssueDate       materials.TextField
	DueDate         materials.TextField

	dueDateOverride      bool
	dueDatePreviousValue string

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (f *UtilitiesInvoiceForm) Load(
	invoice avisha.UtilityInvoice,
	// previousReading will be used to calculate consumption as a subtraction
	// from the current reading.
	previousReading int,
) {
	f.Invoice = invoice
	f.PreviousReading.SetText(strconv.Itoa(previousReading))
	f.Form.Load([]widget.Field{
		{
			Value: widget.IntValuer{
				Value:   &f.Invoice.Reading,
				Default: previousReading,
			},
			Input: &f.CurrentReading,
		},
		{
			Value: widget.CurrencyValuer{
				Value:   &f.Invoice.UnitCost,
				Default: currency.Dollar * 1,
			},
			Input: &f.UnitCost,
		},
		{
			Value: widget.DateValuer{
				Value:   &f.Invoice.Issued,
				Default: time.Now(),
			},
			Input: &f.IssueDate,
		},
		{
			// @Todo pull net from config?
			Value: widget.DateValuer{
				Value: &f.Invoice.Due,
			},
			Input: &f.DueDate,
		},
		{
			Value: widget.IntValuer{
				Value: &f.Invoice.UnitsConsumed,
			},
			Input: &f.UnitsConsumed,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Invoice.Bill,
			},
			Input: &f.TotalCost,
		},
	})
}

func (f *UtilitiesInvoiceForm) Submit() (invoice avisha.UtilityInvoice, ok bool) {
	return f.Invoice, f.Form.Submit()
}

func (f *UtilitiesInvoiceForm) Clear() {
	f.Form.Clear()
}

func (f *UtilitiesInvoiceForm) Update(gtx C) {
	// Compute DueDate unless manually overridden.
	{
		if f.DueDate.Focused() && f.dueDatePreviousValue != f.DueDate.Text() {
			f.dueDateOverride = true
		}
		if !f.dueDateOverride {
			f.DueDate.SetText(func() string {
				issued, err := util.ParseDate(f.IssueDate.Text())
				if err != nil {
					return f.DueDate.Text()
				}
				return util.FormatTime(issued.Add(time.Hour * 24 * 14))
			}())
			f.dueDatePreviousValue = f.DueDate.Text()
		}
	}
	f.UnitsConsumed.SetText(func() string {
		current, err := util.ParseInt(f.CurrentReading.Text())
		if err != nil {
			return "0"
		}
		previous, err := util.ParseInt(f.PreviousReading.Text())
		if err != nil {
			return "0"
		}
		return strconv.Itoa(current - previous)
	}())
	f.TotalCost.SetText(func() string {
		consumed, err := util.ParseInt(f.UnitsConsumed.Text())
		if err != nil {
			return "0"
		}
		cost, err := util.ParseCurrency(f.UnitCost.Text())
		if err != nil {
			return "0"
		}
		total := cost * currency.Currency(consumed)
		return strconv.FormatFloat(total.Dollars(), 'f', 2, 64)
	}())
	f.Form.Validate(gtx)
}

func (f *UtilitiesInvoiceForm) Layout(gtx C, th *style.Theme) D {
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
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Flexed(1, func(gtx C) D {
							return f.IssueDate.Layout(gtx, th.Dark(), "Issue Date")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							return f.DueDate.Layout(gtx, th.Dark(), "Due Date (net 14)")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					f.UnitCost.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return f.UnitCost.Layout(gtx, th.Dark(), "Unit Cost")
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Flexed(1, func(gtx C) D {
							gtx.Queue = nil
							return f.PreviousReading.Layout(gtx, th.Dark(), "Previous Reading")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							return f.CurrentReading.Layout(gtx, th.Dark(), "Current Reading")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Flexed(1, func(gtx C) D {
							gtx.Queue = nil
							return f.UnitsConsumed.Layout(gtx, th.Dark(), "Units Consumed")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							gtx.Queue = nil
							f.TotalCost.Prefix = func(gtx C) D {
								return material.Body1(th.Dark(), "$").Layout(gtx)
							}
							return f.TotalCost.Layout(gtx, th.Dark(), "Total Cost")
						}),
					)
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
