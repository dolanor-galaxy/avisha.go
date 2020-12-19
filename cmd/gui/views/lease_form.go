package views

import (
	"fmt"
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
)

// LeaseForm performs data mutations on a Lease entity.
type LeaseForm struct {
	// Entity data.
	lease  *avisha.Lease
	site   avisha.Site
	tenant avisha.Tenant

	// Finder functions that can search for entities.
	// Allows realtime validation of entities.
	TenantFinder func(name string) (t avisha.Tenant, exists bool)
	SiteFinder   func(number string) (s avisha.Site, exists bool)

	// Form fields.
	Tenant materials.TextField
	Site   materials.TextField
	Date   materials.TextField
	Days   materials.TextField
	Rent   materials.TextField

	// Actions.
	SubmitBtn widget.Clickable
	CancelBtn widget.Clickable
}

// Clear resets the form to previous values.
// If the form is in create mode (no lease entity attached), it empties the
// fields.
// @Todo make clear-reset pattern common accross forms.
func (l *LeaseForm) Clear() {
	if l.lease == nil {
		l.Tenant.Clear()
		l.Site.Clear()
		l.Days.Clear()
		l.Rent.Clear()
		l.Date.Clear()
	} else {
		l.Load(l.lease, l.tenant, l.site)
	}
}

// Load form data from a lease entity.
func (l *LeaseForm) Load(lease *avisha.Lease, tenant avisha.Tenant, site avisha.Site) {
	l.lease = lease
	l.tenant = tenant
	l.site = site
	l.Tenant.SetText(l.tenant.Name)
	l.Site.SetText(l.site.Number)
	l.Days.SetText(strconv.Itoa(l.lease.Term.Days))
	l.Rent.SetText(strconv.Itoa(int(l.lease.Rent)))
	start := lease.Term.Start
	l.Date.SetText(fmt.Sprintf("%d/%d/%d", start.Day(), start.Month(), start.Year()))
}

func (l *LeaseForm) Update(gtx C) {
	for range l.Tenant.Events() {
		if _, err := l.validateTenant(); err != nil {
			l.Tenant.SetError(err.Error())
		} else {
			l.Tenant.ClearError()
		}
	}
	for range l.Site.Events() {
		if _, err := l.validateSite(); err != nil {
			l.Site.SetError(err.Error())
		} else {
			l.Site.ClearError()
		}
	}
	for range l.Date.Events() {
		if _, err := l.validateDate(); err != nil {
			l.Date.SetError(err.Error())
		} else {
			l.Date.ClearError()
		}
	}
	for range l.Days.Events() {
		if _, err := l.validateDays(); err != nil {
			l.Days.SetError(err.Error())
		} else {
			l.Days.ClearError()
		}
	}
	for range l.Rent.Events() {
		if _, err := l.validateRent(); err != nil {
			l.Rent.SetError(err.Error())
		} else {
			l.Rent.ClearError()
		}
	}
}

func (l *LeaseForm) Layout(gtx C, th *style.Theme) D {
	l.Update(gtx)
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
					return l.Days.Layout(gtx, th.Dark(), "Duration (days)")
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
							if l.lease == nil {
								return material.Button(th.Success(), &l.SubmitBtn, "Create").Layout(gtx)
							}
							return material.Button(th.Primary(), &l.SubmitBtn, "Update").Layout(gtx)
						}),
					)
				})
		}),
	)
}

// Submit validates the input data and returns a boolean indicating validity.
func (l *LeaseForm) Submit() (lease avisha.Lease, ok bool) {
	ok = true
	if l.lease != nil {
		lease.ID = l.lease.ID
	}
	if t, err := l.validateTenant(); err != nil {
		l.Tenant.SetError(err.Error())
		ok = false
	} else {
		lease.Tenant = t.ID
	}
	if s, err := l.validateSite(); err != nil {
		l.Site.SetError(err.Error())
		ok = false
	} else {
		lease.Site = s.ID
	}
	if date, err := l.validateDate(); err != nil {
		l.Date.SetError(err.Error())
		ok = false
	} else {
		lease.Term.Start = date
	}
	if days, err := l.validateDays(); err != nil {
		l.Days.SetError(err.Error())
		ok = false
	} else {
		lease.Term.Days = days
	}
	if rent, err := l.validateRent(); err != nil {
		l.Rent.SetError(err.Error())
		ok = false
	} else {
		lease.Rent = rent
	}
	return lease, ok
}

func (l *LeaseForm) validateTenant() (t avisha.Tenant, err error) {
	if find := l.TenantFinder; find != nil {
		if t, ok := find(l.Tenant.Text()); ok {
			return t, nil
		} else {
			return t, fmt.Errorf("not found")
		}
	}
	return t, nil
}

func (l *LeaseForm) validateSite() (s avisha.Site, err error) {
	if find := l.SiteFinder; find != nil {
		if s, ok := find(l.Site.Text()); ok {
			return s, nil
		} else {
			return s, fmt.Errorf("not found")
		}
	}
	return s, nil
}

func (l *LeaseForm) validateDate() (date time.Time, err error) {
	return util.ParseDate(l.Date.Text())
}

func (l *LeaseForm) validateDays() (days int, err error) {
	n, err := strconv.Atoi(l.Days.Text())
	if err != nil {
		return days, fmt.Errorf("days must be a number")
	}
	return n, nil
}

func (l *LeaseForm) validateRent() (rent uint, err error) {
	n, err := strconv.Atoi(l.Rent.Text())
	if err != nil {
		return rent, fmt.Errorf("rent must be a number")
	}
	if n < 0 {
		return rent, fmt.Errorf("must be a positive number")
	}
	return uint(n), nil
}
