package views

// @Todo Form field abstraction; consider if we can abstract forms in a way that
// reduces the per-field boilerplate.
// 1. Realtime validation (run per event, per field)
// 2. Validation function (pure validation logic, common funcs for "isnumber", "isemail", "isdate", etc)
// 3. Submission validation (validate all fields and collect errors)

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
	UnitsConsumed   materials.TextField
	PreviousReading materials.TextField
	CurrentReading  materials.TextField
	UnitCost        materials.TextField
	TotalCost       materials.TextField
	IssueDate       materials.TextField
	DueDate         materials.TextField

	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
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

func (f *UtilitiesInvoiceForm) Clear() {
	f.IssueDate.Clear()
	f.IssueDate.SetText(util.FormatTime(time.Now()))
	f.DueDate.Clear()
	f.DueDate.SetText(util.FormatTime(time.Now().Add(time.Hour * 24 * 14)))
	// @Todo load unit cost default value from storage.
	f.UnitCost.Clear()
	f.UnitCost.SetText("1")
	f.UnitsConsumed.Clear()
}

func (f *UtilitiesInvoiceForm) Update(gtx C) {
	for range f.UnitCost.Events() {
		if _, err := f.validateUnitCost(); err != nil {
			f.UnitCost.SetError(err.Error())
		} else {
			f.UnitCost.ClearError()
		}
	}
	for range f.IssueDate.Events() {
		if t, err := f.validateIssueDate(); err != nil {
			f.IssueDate.SetError(err.Error())
		} else {
			f.IssueDate.ClearError()
			f.DueDate.SetText(util.FormatTime(t.Add(time.Hour * 24 * 14)))
		}
	}
	for range f.DueDate.Events() {
		if _, err := f.validateDueDate(); err != nil {
			f.DueDate.SetError(err.Error())
		} else {
			f.DueDate.ClearError()
		}
	}
	for range f.CurrentReading.Events() {
		if current, err := f.validateCurrentReading(); err != nil {
			f.CurrentReading.SetError(err.Error())
		} else {
			f.CurrentReading.ClearError()
			if previous, err := f.validatePreviousReading(); err != nil {
				f.PreviousReading.SetError(err.Error())
			} else {
				f.PreviousReading.ClearError()
				consumed := current - previous
				f.UnitsConsumed.SetText(strconv.Itoa(consumed))
				if cost, err := f.validateUnitCost(); err != nil {
					f.UnitCost.SetError(err.Error())
				} else {
					f.UnitCost.ClearError()
					f.TotalCost.SetText(strconv.Itoa(consumed * cost))
				}
			}
		}
	}
}

func (f *UtilitiesInvoiceForm) Submit() (invoice avisha.UtilityInvoice, ok bool) {
	ok = true
	if unitCost, err := f.validateUnitCost(); err != nil {
		f.UnitCost.SetError(err.Error())
		ok = false
	} else {
		invoice.UnitCost = currency.Dollar * currency.Currency(unitCost)
	}
	if unitsConsumed, err := f.validateUnitsConsumed(); err != nil {
		f.UnitsConsumed.SetError(err.Error())
		ok = false
	} else {
		invoice.UnitsConsumed = unitsConsumed
	}
	if issueDate, err := f.validateIssueDate(); err != nil {
		f.UnitsConsumed.SetError(err.Error())
		ok = false
	} else {
		invoice.Issued = issueDate
	}
	if dueDate, err := f.validateDueDate(); err != nil {
		f.UnitsConsumed.SetError(err.Error())
		ok = false
	} else {
		invoice.Due = dueDate
	}
	invoice.Bill = (invoice.UnitCost * currency.Currency(invoice.UnitsConsumed))
	return invoice, ok
}

func (f *UtilitiesInvoiceForm) validateUnitCost() (int, error) {
	return util.ParseInt(f.UnitCost.Text())
}

func (f *UtilitiesInvoiceForm) validateUnitsConsumed() (int, error) {
	return util.ParseInt(f.UnitsConsumed.Text())
}

func (f *UtilitiesInvoiceForm) validatePreviousReading() (int, error) {
	return util.ParseInt(f.PreviousReading.Text())
}

func (f *UtilitiesInvoiceForm) validateCurrentReading() (int, error) {
	return util.ParseInt(f.CurrentReading.Text())
}

func (f *UtilitiesInvoiceForm) validateIssueDate() (time.Time, error) {
	return util.ParseDate(f.IssueDate.Text())
}

func (f *UtilitiesInvoiceForm) validateDueDate() (time.Time, error) {
	return util.ParseDate(f.DueDate.Text())
}
