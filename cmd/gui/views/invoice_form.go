package views

import (
	"fmt"
	"image"
	"strconv"
	"strings"
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

	// Bill is the final amount due.
	Bill materials.TextField

	// Charges.
	LateFee    materials.TextField
	LineCharge materials.TextField
	GST        materials.TextField
	Activity   materials.TextField

	IssueDate materials.TextField
	DueDate   materials.TextField

	dueDateOverride      bool
	dueDatePreviousValue string

	invoiceNet time.Duration

	Form      widget.Form
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

func (f *UtilitiesInvoiceForm) Load(
	invoice avisha.UtilityInvoice,
	settings avisha.Settings,
	previousReading int,
) {
	f.Invoice = invoice
	f.invoiceNet = settings.Defaults.InvoiceNet
	f.Invoice.GST = settings.Defaults.GST
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
				Value: &f.Invoice.Charges.LateFee,
			},
			Input: &f.LateFee,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Invoice.Charges.LineCharge,
			},
			Input: &f.LineCharge,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Invoice.Charges.GST,
			},
			Input: &f.GST,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Invoice.Charges.Activity,
			},
			Input: &f.Activity,
		},
		{
			Value: widget.CurrencyValuer{
				Value: &f.Invoice.Bill,
			},
			Input: &f.Bill,
		},
	})
}

func (f *UtilitiesInvoiceForm) Submit() (invoice avisha.UtilityInvoice, ok bool) {
	f.Invoice.Balance.Debit(f.Invoice.Bill)
	return f.Invoice, f.Form.Submit()
}

func (f *UtilitiesInvoiceForm) Clear() {
	f.Form.Clear()
	f.dueDateOverride = false
	f.invoiceNet = 0
	f.dueDatePreviousValue = ""
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
				return util.FormatTime(issued.Add(f.invoiceNet))
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
	f.Activity.SetText(func() string {
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
	f.GST.SetText(func() string {
		activity, err := util.ParseCurrency(f.Activity.Text())
		if err != nil {
			return "0"
		}
		latefee, err := util.ParseCurrency(f.LateFee.Text())
		if err != nil {
			return "0"
		}
		linecharge, err := util.ParseCurrency(f.LineCharge.Text())
		if err != nil {
			return "0"
		}
		total := activity + latefee + linecharge
		return strings.TrimPrefix(currency.Currency(float64(total)*(f.Invoice.GST/100)).String(), "$")
	}())
	f.Bill.SetText(func() string {
		activity, err := util.ParseCurrency(f.Activity.Text())
		if err != nil {
			return "0"
		}
		latefee, err := util.ParseCurrency(f.LateFee.Text())
		if err != nil {
			return "0"
		}
		linecharge, err := util.ParseCurrency(f.LineCharge.Text())
		if err != nil {
			return "0"
		}
		gst, err := util.ParseCurrency(f.GST.Text())
		if err != nil {
			return "0"
		}
		total := activity + latefee + linecharge + gst
		return strings.TrimPrefix(total.String(), "$")
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
							return f.DueDate.Layout(gtx, th.Dark(), fmt.Sprintf("Due Date (net %d)", int(f.invoiceNet.Hours()/24)))
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
							f.Activity.Prefix = func(gtx C) D {
								return material.Body1(th.Dark(), "$").Layout(gtx)
							}
							return f.Activity.Layout(gtx, th.Dark(), "Activity")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					f.LineCharge.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return f.LineCharge.Layout(gtx, th.Dark(), "LineCharge")
				}),
				layout.Rigid(func(gtx C) D {
					f.LateFee.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return f.LateFee.Layout(gtx, th.Dark(), "LateFee")
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Queue = nil
					f.GST.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return f.GST.Layout(gtx, th.Dark(), fmt.Sprintf("GST (%.0f%%)", f.Invoice.GST))
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Queue = nil
					f.Bill.Prefix = func(gtx C) D {
						return material.Body1(th.Dark(), "$").Layout(gtx)
					}
					return f.Bill.Layout(gtx, th.Dark(), "Bill")
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
