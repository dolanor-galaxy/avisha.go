package views

import (
	"fmt"
	"image"
	"log"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
)

// LeaseForm performs data mutations on a Lease entity.
type LeaseForm struct {
	nav.Route
	App *avisha.App
	Th  *style.Theme

	lease  *avisha.Lease
	site   avisha.Site
	tenant avisha.Tenant

	Tenant materials.TextField
	Site   materials.TextField
	Date   materials.TextField
	Days   materials.TextField
	Rent   materials.TextField
	Submit widget.Clickable
	Cancel widget.Clickable
}

func (l *LeaseForm) Title() string {
	return "Lease Form"
}

// TODO: route back on error?
func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		l.lease = lease
		if err := l.App.One("ID", lease.Tenant, &l.tenant); err != nil {
			log.Printf("loading tenant: %+v: %v", lease.Tenant, err)
			return
		}
		if err := l.App.One("ID", lease.Site, &l.site); err != nil {
			log.Printf("loading site: %v", err)
			return
		}
		l.Tenant.SetText(l.tenant.Name)
		l.Site.SetText(l.site.Number)
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
					label := material.Label(l.Th.Primary(), unit.Dp(24), fmt.Sprintf("%s-%s", l.tenant.Name, l.site.Number))
					label.Alignment = text.Middle
					label.Color = l.Th.Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (l *LeaseForm) Update(gtx C) {
	clear := func() {
		l.Tenant.SetText("")
		l.Site.SetText("")
		l.Days.SetText("")
		l.Rent.SetText("")
		l.Date.SetText("")
		l.lease = nil
	}
	if l.Submit.Clicked() {
		if err := l.submit(); err != nil {
			log.Printf("submitting lease: %v", err)
		} else {
			clear()
			l.Route.Back()
		}
	}
	if l.Cancel.Clicked() {
		clear()
		l.Route.Back()
	}
}

func (l *LeaseForm) Layout(gtx C) D {
	l.Update(gtx)
	l.Tenant.Validator = func(text string) string {
		var t avisha.Tenant
		if err := l.App.One("Name", text, &t); err != nil {
			return err.Error()
		}
		return ""
	}
	l.Site.Validator = func(text string) string {
		var s avisha.Site
		if err := l.App.One("Number", text, &s); err != nil {
			return err.Error()
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
							return l.Tenant.Layout(gtx, l.Th.Primary(), "Tenant")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							if l.lease != nil {
								gtx.Queue = nil
							}
							return l.Site.Layout(gtx, l.Th.Primary(), "Site")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if l.lease != nil {
						gtx.Queue = nil
					}
					return l.Date.Layout(gtx, l.Th.Primary(), "Start Date")
				}),
				layout.Rigid(func(gtx C) D {
					if l.lease != nil {
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
							if l.lease == nil {
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
//
// FIXME: figure out why errors wont display when submitting an empty form.
func (l *LeaseForm) validate() (*avisha.Lease, bool) {
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
		l.Days.SetError("rent not a number")
	}
	lease.Tenant, ok = func() (id int, ok bool) {
		var t avisha.Tenant
		if l.Tenant.Text() == "" {
			l.Tenant.SetError("tenant cannot be empty")
			return t.ID, false
		}
		err := l.App.One("Name", l.Tenant.Text(), &t)
		if err != nil {
			log.Printf("finding tenant for lease: %v", err)
		}
		return t.ID, err == nil
	}()
	if ok {
		invalid = false
	}
	lease.Site, ok = func() (site int, ok bool) {
		var s avisha.Site
		if l.Site.Text() == "" {
			l.Tenant.SetError("site cannot be empty")
			return s.ID, false
		}
		err := l.App.One("Number", l.Site.Text(), &s)
		if err != nil {
			log.Printf("finding site for lease: %v", err)
		}
		return s.ID, err == nil
	}()
	if ok {
		invalid = false
	}
	return lease, !invalid && err == nil
}

func (l *LeaseForm) submit() error {
	lease, ok := l.validate()
	if !ok {
		return fmt.Errorf("failed validation")
	}
	if l.lease == nil {
		if err := l.App.CreateLease(lease); err != nil {
			return fmt.Errorf("creating lease: %v", err)
		}
	} else {
		if err := l.App.Update(lease); err != nil {
			return fmt.Errorf("updating lease: %v", err)
		}
	}
	return nil
}
