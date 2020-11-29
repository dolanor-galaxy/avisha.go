package views

import (
	"fmt"
	"image"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
	"github.com/jackmordaunt/avisha-fn/storage"
)

// LeaseForm performs data mutations on a Lease entity.
type LeaseForm struct {
	nav.Route
	App   *avisha.App
	Th    *style.Theme
	lease *avisha.Lease

	Tenant materials.TextField
	Site   materials.TextField
	Date   materials.TextField
	Days   materials.TextField
	Rent   materials.TextField
	Submit widget.Clickable
	Cancel widget.Clickable

	// creating if in "create mode", which means the entity doesn't exist.
	// This determines whether submissions creates or updates an entity.
	creating bool
}

func (l *LeaseForm) Title() string {
	return "Lease Form"
}

func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		l.lease = lease
		l.creating = lease.Cmp() == (avisha.Lease{}).Cmp()
		l.Tenant.SetText(l.lease.Tenant)
		l.Site.SetText(l.lease.Site)
		l.Days.SetText(strconv.Itoa(l.lease.Term.Days))
		l.Rent.SetText(strconv.Itoa(int(l.lease.Rent)))
		start := lease.Term.Start
		l.Date.SetText(fmt.Sprintf("%d/%d/%d", start.Day(), start.Month(), start.Year()))
	}
}

func (l *LeaseForm) Context() (list []layout.Widget) {
	if l.lease != nil {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(l.Th.Primary(), unit.Dp(24), l.lease.ID())
					label.Alignment = text.Middle
					label.Color = l.Th.Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (l *LeaseForm) Update(gtx C) {
	if l.Submit.Clicked() {
		l.submit()
	}
	if l.Cancel.Clicked() {
		l.Route.Back()
	}
}

func (l *LeaseForm) Layout(gtx C) D {
	l.Update(gtx)
	l.Tenant.Validator = func(text string) string {
		if _, ok := l.App.Query(func(ent storage.Entity) bool {
			if tenant, ok := ent.(*avisha.Tenant); ok {
				return tenant.Name == text
			}
			return false
		}); !ok {
			return fmt.Sprintf("tenant %q does not exist", text)
		}
		return ""
	}
	l.Site.Validator = func(text string) string {
		if _, ok := l.App.Query(func(ent storage.Entity) bool {
			if site, ok := ent.(*avisha.Site); ok {
				return site.Number == text
			}
			return false
		}); !ok {
			return fmt.Sprintf("site %q does not exist", text)
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
							if !l.creating {
								gtx.Queue = nil
							}
							return l.Tenant.Layout(gtx, l.Th.Primary(), "Tenant")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							if !l.creating {
								gtx.Queue = nil
							}
							return l.Site.Layout(gtx, l.Th.Primary(), "Site")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if !l.creating {
						gtx.Queue = nil
					}
					return l.Date.Layout(gtx, l.Th.Primary(), "Start Date")
				}),
				layout.Rigid(func(gtx C) D {
					if !l.creating {
						gtx.Queue = nil
					}
					return l.Days.Layout(gtx, l.Th.Primary(), "Duration (days)")
				}),
				layout.Rigid(func(gtx C) D {
					l.Rent.Prefix = func(gtx C) D {
						return material.Body1(l.Th.Primary(), "$").Layout(gtx)
					}
					return l.Rent.Layout(gtx, l.Th.Primary(), "Rent (weekly)")
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
							return material.Button(l.Th.Secondary(), &l.Cancel, "Cancel").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							if l.creating {
								return material.Button(l.Th.Success(), &l.Submit, "Create").Layout(gtx)
							}
							return material.Button(l.Th.Primary(), &l.Submit, "Update").Layout(gtx)
						}),
					)
				})
		}),
	)
}

// validate form date returning true if the form is ready to be submitted.
// Data is saved to the embedded entity.
func (l *LeaseForm) validate() bool {
	var err error
	if l.Tenant.Text() == "" {
		l.Tenant.SetError("tenant cannot be empty")
		return false
	}
	if l.Site.Text() == "" {
		l.Tenant.SetError("site cannot be empty")
		return false
	}
	l.lease.Term.Start, err = func() (time.Time, error) {
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
		return false
	}
	l.lease.Term.Days, err = strconv.Atoi(l.Days.Text())
	if err != nil {
		l.Days.SetError("days not a number")
		return false
	}
	l.lease.Rent, err = func() (uint, error) {
		n, err := strconv.Atoi(l.Rent.Text())
		return uint(n), err
	}()
	if err != nil {
		l.Days.SetError("rent not a number")
		return false
	}
	return true
}

func (l *LeaseForm) submit() {
	if isValid := l.validate(); !isValid {
		return
	}
	if l.creating {
		if err := l.App.CreateLease(
			l.Tenant.Text(),
			l.Site.Text(),
			l.lease.Term,
			l.lease.Rent,
		); err != nil {
			fmt.Printf("creating lease: %v\n", err)
			return
		}
		l.Route.Back()
	} else {
		if err := l.App.ChangeRent(
			l.Tenant.Text(),
			l.Site.Text(),
			l.lease.Term,
			l.lease.Rent,
		); err != nil {
			fmt.Printf("changing rent: %v\n", err)
			return
		}
		l.Route.Back()
	}
}
