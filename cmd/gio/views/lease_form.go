package views

import (
	"fmt"
	"image"
	"log"
	"strconv"
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
	Date   DateInput
	Days   materials.TextField
	Rent   materials.TextField
	Submit widget.Clickable
	Cancel widget.Clickable
}

// DateInput allows inputting of date information textually.
// Composed of three inputs: day, month, year.
type DateInput struct {
	Day   materials.TextField
	Month materials.TextField
	Year  materials.TextField
}

func (input *DateInput) Set(date time.Time) {
	var (
		day, month, year = "", "", ""
	)
	if date != (time.Time{}) {
		day = fmt.Sprintf("%d", date.Day())
		month = fmt.Sprintf("%d", date.Month())
		year = fmt.Sprintf("%d", date.Year())
	}
	input.Day.SetText(day)
	input.Month.SetText(month)
	input.Year.SetText(year)
}

func (input *DateInput) Date() (time.Time, error) {
	year, err := strconv.Atoi(input.Year.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("year: not a number")
	}
	if year < 0 {
		return time.Time{}, fmt.Errorf("year: out of bounds (must be positive number) got %d", year)
	}
	month, err := strconv.Atoi(input.Month.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("month: not a number")
	}
	if month < int(time.January) || month > int(time.December) {
		return time.Time{}, fmt.Errorf("month: out of bounds (1-12) got %d", month)
	}
	day, err := strconv.Atoi(input.Day.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("day: not a number")
	}
	if day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("day: out of bounds (1-31) got %d", day)
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

func (input *DateInput) Layout(gtx C, th *style.Theme) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(
		gtx,
		layout.Flexed(1, func(gtx C) D {
			return input.Day.Layout(gtx, th.Primary(), "Day")
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
		}),
		layout.Flexed(1, func(gtx C) D {
			return input.Month.Layout(gtx, th.Primary(), "Month")
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
		}),
		layout.Flexed(1, func(gtx C) D {
			return input.Year.Layout(gtx, th.Primary(), "Year")
		}),
	)
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
		l.Date.Set(lease.Term.Start)
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
					return l.Tenant.Layout(gtx, l.Th.Primary(), "Tenant")
				}),
				layout.Rigid(func(gtx C) D {
					return l.Site.Layout(gtx, l.Th.Primary(), "Site")
				}),
				layout.Rigid(func(gtx C) D {
					return l.Date.Layout(gtx, l.Th)
				}),
				layout.Rigid(func(gtx C) D {
					return l.Days.Layout(gtx, l.Th.Primary(), "Days")
				}),
				layout.Rigid(func(gtx C) D {
					return l.Rent.Layout(gtx, l.Th.Primary(), "Rent")
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
	start, err := l.Date.Date()
	if err != nil {
		return fmt.Errorf("start date: %w", err)
	}
	days, err := strconv.Atoi(l.Days.Text())
	if err != nil {
		return fmt.Errorf("days not a number: %w", err)
	}
	rent, err := strconv.Atoi(l.Rent.Text())
	if err != nil {
		return fmt.Errorf("rent not a number: %w", err)
	}
	if err := l.App.CreateLease(
		l.Tenant.Text(),
		l.Site.Text(),
		avisha.Term{Start: start, Days: days},
		uint(rent),
	); err != nil {
		return fmt.Errorf("creating lease: %w", err)
	}
	return nil
}
