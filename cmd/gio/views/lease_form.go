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
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
)

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
}

func (l *LeaseForm) Title() string {
	return "Lease Form"
}

func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		l.lease = lease
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

// Creating reports whether the form is creating a new entity.
// Returns false when the entity already exists.
func (l *LeaseForm) Creating() bool {
	return l.lease.Cmp() == (avisha.Lease{}).Cmp()
}

func (l *LeaseForm) Update(gtx C) {
	if l.Submit.Clicked() {
		if err := l.submit(); err != nil {
			// give error to app or render under field.
			log.Printf("submitting lease form: %v", err)
		}
		l.Route.Back()
	}
	if l.Cancel.Clicked() {
		l.Route.Back()
	}
}

func (l *LeaseForm) Layout(gtx C) D {
	l.Update(gtx)
	// TODO: implement disabled text field states.
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
							return l.Tenant.Layout(gtx, l.Th.Primary(), "Tenant")
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
						}),
						layout.Flexed(1, func(gtx C) D {
							return l.Site.Layout(gtx, l.Th.Primary(), "Site")
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return l.Date.Layout(gtx, l.Th.Primary(), "Start Date")
				}),
				layout.Rigid(func(gtx C) D {
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
							return material.Button(l.Th.Primary(), &l.Submit, "Submit").Layout(gtx)
						}),
					)
				})
		}),
	)
}

func (l *LeaseForm) submit() error {
	s := l.Date.Text()
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return fmt.Errorf("start date: invalid format: must be dd/mm/yyyy")
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("start date: year not a number: %s", parts[2])
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("start date: month not a number: %s", parts[2])
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("start date: day not a number: %s", parts[2])
	}
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	days, err := strconv.Atoi(l.Days.Text())
	if err != nil {
		return fmt.Errorf("days not a number: %w", err)
	}
	rent, err := strconv.Atoi(l.Rent.Text())
	if err != nil {
		return fmt.Errorf("rent not a number: %w", err)
	}
	if l.Creating() {
		if err := l.App.CreateLease(
			l.Tenant.Text(),
			l.Site.Text(),
			avisha.Term{Start: start, Days: days},
			uint(rent),
		); err != nil {
			return fmt.Errorf("creating lease: %w", err)
		}
	} else {
		if err := l.App.ChangeRent(
			l.Tenant.Text(),
			l.Site.Text(),
			avisha.Term{Start: start, Days: days},
			uint(rent),
		); err != nil {
			return fmt.Errorf("changing rent: %w", err)
		}
	}
	return nil
}
