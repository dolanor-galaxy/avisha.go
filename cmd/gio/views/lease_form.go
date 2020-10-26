package views

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

type LeaseForm struct {
	*avisha.App
	*material.Theme
	lease *avisha.Lease

	Tenant materials.TextField
	Site   materials.TextField

	Start materials.TextField
	Days  materials.TextField

	Rent materials.TextField

	Submit widget.Clickable
	Cancel widget.Clickable

	layout.List

	reroute string
}

func (l *LeaseForm) ReRoute() (string, interface{}) {
	defer func() { l.reroute = "" }()
	return l.reroute, nil
}

func (l *LeaseForm) Receive(data interface{}) {
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		l.lease = lease
		l.Tenant.SetText(l.lease.Tenant)
		l.Site.SetText(l.lease.Site)
		l.Start.SetText(l.lease.Term.Start.String())
		l.Days.SetText(strconv.Itoa(l.lease.Term.Days))
		l.Rent.SetText(strconv.Itoa(int(l.lease.Rent)))
	}
}

func (l *LeaseForm) Context() (list []layout.Widget) {
	if l.lease != nil {
		list = append(list, func(gtx Ctx) Dims {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx Ctx) Dims {
					label := material.Label(l.Theme, unit.Dp(24), l.lease.ID())
					label.Alignment = text.Middle
					label.Color = l.Theme.Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (l *LeaseForm) Update(gtx Ctx) {
	if l.Submit.Clicked() {
		if err := l.submit(); err != nil {
			// give error to app or render under field.
			log.Printf("submitting lease form: %v", err)
		}
	}
	if l.Cancel.Clicked() {
		l.reroute = "back"
	}
}

func (l *LeaseForm) Layout(gtx Ctx) Dims {
	l.Update(gtx)
	l.List.Axis = layout.Vertical
	return layout.UniformInset(unit.Dp(10)).Layout(
		gtx,
		func(gtx Ctx) Dims {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(
				gtx,
				layout.Flexed(1, func(gtx Ctx) Dims {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx Ctx) Dims {
							return l.Tenant.Layout(gtx, l.Theme, "Tenant")
						}),
						layout.Rigid(func(gtx Ctx) Dims {
							return l.Site.Layout(gtx, l.Theme, "Site")
						}),
						layout.Rigid(func(gtx Ctx) Dims {
							return l.Start.Layout(gtx, l.Theme, "Start")
						}),
						layout.Rigid(func(gtx Ctx) Dims {
							return l.Days.Layout(gtx, l.Theme, "Days")
						}),
						layout.Rigid(func(gtx Ctx) Dims {
							return l.Rent.Layout(gtx, l.Theme, "Rent")
						}),
					)
				}),
				layout.Rigid(func(gtx Ctx) Dims {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx Ctx) Dims {
							return material.Button(l.Theme, &l.Cancel, "Cancel").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx Ctx) Dims {
							return Dims{Size: gtx.Constraints.Min}
						}),
						layout.Rigid(func(gtx Ctx) Dims {
							return material.Button(l.Theme, &l.Submit, "Submit").Layout(gtx)
						}),
					)
				}),
			)
		},
	)
}

func (l *LeaseForm) submit() error {
	start, err := time.Parse("", l.Start.Text())
	if err != nil {
		return fmt.Errorf("invalid date specifier: %w", err)
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
