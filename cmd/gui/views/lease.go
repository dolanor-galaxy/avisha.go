package views

import (
	"fmt"
	"log"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
)

// LeasePage contains actions for interacting with a lease including data entry
// and service payments.
type LeasePage struct {
	nav.Route
	App  *avisha.App
	Th   *style.Theme
	Form LeaseForm
}

func (page *LeasePage) Title() string {
	return "Lease"
}

// TODO: route back on error?
func (p *LeasePage) Receive(data interface{}) {
	p.Form.TenantFinder = func(name string) (t avisha.Tenant, ok bool) {
		err := p.App.One("Name", name, &t)
		return t, err == nil
	}
	p.Form.SiteFinder = func(number string) (s avisha.Site, ok bool) {
		err := p.App.One("Number", number, &s)
		return s, err == nil
	}
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		p.Form.lease = lease
		if err := p.App.One("ID", lease.Tenant, &p.Form.tenant); err != nil {
			log.Printf("loading tenant: %+v: %v", lease.Tenant, err)
			return
		}
		if err := p.App.One("ID", lease.Site, &p.Form.site); err != nil {
			log.Printf("loading site: %v", err)
			return
		}
		p.Form.Tenant.SetText(p.Form.tenant.Name)
		p.Form.Site.SetText(p.Form.site.Number)
		p.Form.Days.SetText(strconv.Itoa(p.Form.lease.Term.Days))
		p.Form.Rent.SetText(strconv.Itoa(int(p.Form.lease.Rent)))
		start := lease.Term.Start
		p.Form.Date.SetText(fmt.Sprintf("%d/%d/%d", start.Day(), start.Month(), start.Year()))
	}
}

func (p *LeasePage) Context() (list []layout.Widget) {
	if p.Form.lease != nil {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(
						p.Th.Primary(),
						unit.Dp(24),
						fmt.Sprintf("%s-%s", p.Form.tenant.Name, p.Form.site.Number))
					label.Alignment = text.Middle
					label.Color = p.Th.Color.InvText
					return label.Layout(gtx)
				})
		})
	}
	return list
}

func (p *LeasePage) Update(gtx C) {
	if p.Form.SubmitBtn.Clicked() {
		if lease, ok := p.Form.Submit(); ok {
			if err := func() error {
				if create := lease.ID == 0; create {
					if err := p.App.CreateLease(&lease); err != nil {
						return fmt.Errorf("creating lease: %w", err)
					}
				} else {
					if err := p.App.Update(lease); err != nil {
						return fmt.Errorf("updating lease: %w", err)
					}
				}
				return nil
			}(); err != nil {
				log.Printf("%v", err)
			} else {
				p.Form.Clear()
				p.Route.Back()
			}
		}
	}
	if p.Form.CancelBtn.Clicked() {
		p.Form.Clear()
		p.Route.Back()
	}
}

func (p *LeasePage) Layout(gtx C) D {
	p.Update(gtx)
	return p.Form.Layout(gtx, p.Th)
}
