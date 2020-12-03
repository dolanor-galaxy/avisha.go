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
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
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

// Clear the form fields.
func (l *LeaseForm) Clear() {
	l.Tenant.Clear()
	l.Site.Clear()
	l.Days.Clear()
	l.Rent.Clear()
	l.Date.Clear()
	l.lease = nil
}

func (l *LeaseForm) Layout(gtx C, th *style.Theme) D {
	l.Tenant.Validator = func(text string) string {
		if l.TenantFinder == nil {
			return ""
		}
		if _, exists := l.TenantFinder(text); !exists {
			return "not found"
		}
		return ""
	}
	l.Site.Validator = func(text string) string {
		if l.SiteFinder == nil {
			return ""
		}
		if _, exists := l.SiteFinder(text); !exists {
			return "not found"
		}
		return ""
	}
	l.Date.Validator = func(text string) string {
		parts := strings.Split(text, "/")
		if len(parts) != 3 {
			return "format must be dd/mm/yyy"
		}
		for ii, part := range parts {
			if _, err := strconv.Atoi(part); err != nil {
				switch ii {
				case 0:
					return "day must be a number"
				case 1:
					return "month must be a number"
				case 2:
					return "year must be a number"
				}
			}
		}
		return ""
	}
	l.Days.Validator = func(text string) string {
		if _, err := strconv.Atoi(text); err != nil {
			return "must be a number"
		}
		return ""
	}
	l.Rent.Validator = func(text string) string {
		if _, err := strconv.Atoi(text); err != nil {
			return "must be a number"
		}
		return ""
	}
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
							if l.lease != nil {
								gtx.Queue = nil
							}
							return l.Tenant.Layout(gtx, th.Primary(), "Tenant")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							if l.lease != nil {
								gtx.Queue = nil
							}
							return l.Site.Layout(gtx, th.Primary(), "Site")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if l.lease != nil {
						gtx.Queue = nil
					}
					return l.Date.Layout(gtx, th.Primary(), "Start Date")
				}),
				layout.Rigid(func(gtx C) D {
					if l.lease != nil {
						gtx.Queue = nil
					}
					return l.Days.Layout(gtx, th.Primary(), "Duration (days)")
				}),
				layout.Rigid(func(gtx C) D {
					l.Rent.Prefix = func(gtx C) D {
						return material.Body1(th.Primary(), "$").Layout(gtx)
					}
					return l.Rent.Layout(gtx, th.Primary(), "Rent (weekly)")
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
//
// FIXME: figure out why errors wont display when submitting an empty form.
func (l *LeaseForm) Submit() (*avisha.Lease, bool) {
	var (
		err     error
		ok      bool
		invalid = true
		lease   = &avisha.Lease{
			ID: func() int {
				if l.lease != nil {
					return l.lease.ID
				}
				return 0
			}(),
		}
	)
	lease.Term.Start, err = func() (time.Time, error) {
		s := l.Date.Text()
		parts := strings.Split(s, "/")
		if len(parts) != 3 {
			return time.Time{}, fmt.Errorf("invalid format: must be dd/mm/yyyy")
		}
		year, err := strconv.Atoi(parts[2])
		if err != nil {
			return time.Time{}, fmt.Errorf("year not a number: %s", parts[2])
		}
		month, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("month not a number: %s", parts[2])
		}
		day, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("day not a number: %s", parts[2])
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
	}()
	if err != nil {
		l.Date.SetError(err.Error())
	}
	lease.Term.Days, err = strconv.Atoi(l.Days.Text())
	if err != nil {
		l.Days.SetError("days not a number")
	}
	lease.Rent, err = func() (uint, error) {
		n, err := strconv.Atoi(l.Rent.Text())
		return uint(n), err
	}()
	if err != nil {
		l.Rent.SetError("rent not a number")
	}
	lease.Tenant, ok = func() (id int, ok bool) {
		if find := l.TenantFinder; find != nil {
			if t, ok := find(l.Tenant.Text()); ok {
				return t.ID, true
			} else {
				l.Tenant.SetError("not found")
			}
		}
		return 0, false
	}()
	if ok {
		invalid = false
	}
	lease.Site, ok = func() (site int, ok bool) {
		if find := l.SiteFinder; find != nil {
			if s, ok := find(l.Site.Text()); ok {
				return s.ID, true
			} else {
				l.Site.SetError("not found")
			}
		}
		return 0, false
	}()
	if ok {
		invalid = false
	}
	return lease, !invalid && err == nil
}
