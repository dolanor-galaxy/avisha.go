package views

import (
	"fmt"
	"image"
	"log"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"
)

// LeasePage contains actions for interacting with a lease including data entry
// and service payments.
type LeasePage struct {
	nav.Route
	App  *avisha.App
	Th   *style.Theme
	Form LeaseForm

	PayUtility  widget.Clickable
	BillUtility widget.Clickable
	PayRent     widget.Clickable
	BillRent    widget.Clickable

	Dialog style.Dialog

	modal layout.Widget
	lease avisha.Lease
}

func (page *LeasePage) Title() string {
	return "Lease"
}

// @Todo: route back on error?
func (p *LeasePage) Receive(data interface{}) {
	p.lease = avisha.Lease{}
	p.Form.TenantFinder = func(name string) (t avisha.Tenant, ok bool) {
		err := p.App.One("Name", name, &t)
		return t, err == nil
	}
	p.Form.SiteFinder = func(number string) (s avisha.Site, ok bool) {
		err := p.App.One("Number", number, &s)
		return s, err == nil
	}
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		// IMPROVEMENT: use one source of lease data.
		// Note: this just copies the data; one copy goes to the form, and one
		// to the page.
		p.lease = *lease
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
	if p.lease.ID != 0 {
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

func (p *LeasePage) Modal(gtx C) D {
	if p.modal == nil {
		return D{}
	}
	return p.modal(gtx)
}

func (p *LeasePage) Update(gtx C) {
	if p.lease.ID != 0 {
		if err := p.App.One("ID", p.lease.ID, &p.lease); err != nil {
			log.Printf("error: loading lease: %d: %v", p.lease.ID, err)
		}
	}
	if p.Form.SubmitBtn.Clicked() {
		if lease, ok := p.Form.Submit(); ok {
			if err := func() error {
				if create := lease.ID == 0; create {
					if err := p.App.CreateLease(&lease); err != nil {
						return fmt.Errorf("creating lease: %w", err)
					}
				} else {
					if err := p.App.Update(&lease); err != nil {
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
	if p.PayUtility.Clicked() {
		p.Dialog.Context = "pay-utilities"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Pay Utilities", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Primary(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	if p.BillUtility.Clicked() {
		p.Dialog.Context = "bill-utilities"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Bill Utilities", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Primary(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	if p.PayRent.Clicked() {
		p.Dialog.Context = "pay-rent"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Pay Rent", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Primary(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	if p.BillRent.Clicked() {
		p.Dialog.Context = "bill-rent"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Bill Rent", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Primary(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	for range p.Dialog.Input.Events() {
		n, err := strconv.Atoi(p.Dialog.Input.Text())
		if err != nil {
			p.Dialog.Input.SetError("must be a valid number")
		} else if n < 1 {
			p.Dialog.Input.SetError("must be an amount greater than 0")
		} else {
			p.Dialog.Input.ClearError()
		}
	}
	if p.Dialog.Ok.Clicked() {
		n, err := func() (int, error) {
			n, err := strconv.Atoi(p.Dialog.Input.Text())
			if err != nil {
				return 0, fmt.Errorf("must be a valid number")
			} else if n < 1 {
				return 0, fmt.Errorf("must be an amount greater than 0")
			}
			return n, nil
		}()
		if err != nil {
			p.Dialog.Input.SetError(err.Error())
		} else {
			p.Dialog.Input.Clear()
			parts := strings.Split(p.Dialog.Context, "-")
			mode, service := parts[0], parts[1]
			switch mode {
			case "pay":
				if err := p.App.PayService(p.lease.ID, service, uint(n)); err != nil {
					log.Printf("paying service: %v", err)
				}
			case "bill":
				if err := p.App.BillService(p.lease.ID, service, uint(n)); err != nil {
					log.Printf("billing service: %v", err)
				}
			}
			p.modal = nil
		}
	}
	if p.Dialog.Cancel.Clicked() {
		p.Dialog.Input.Clear()
		p.modal = nil
	}
	if p.modal != nil {
		// @Improvement: implies that modal must be rendering a dialog; thus any
		// other modal content will call focus on the dialog.
		// 1. does this matter?
		// 2. best way to make this dep explicit?
		p.Dialog.Input.Focus()
	}
}

func (p *LeasePage) Layout(gtx C) D {
	// @Improvement: render different contexts under tabs?
	p.Update(gtx)
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
		return style.Container{
			BreakPoint: unit.Dp(700),
			Constrain:  true,
			// Scroll:     true,
		}.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx C) D {
							if p.lease.ID == 0 {
								return D{}
							}
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(
								gtx,
								layout.Rigid(func(gtx C) D {
									return material.Label(p.Th.Primary(), unit.Dp(20), "Services").Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return D{Size: image.Point{Y: gtx.Px(unit.Dp(10))}}
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
										Spacing:   layout.SpaceBetween,
									}.Layout(
										gtx,
										layout.Flexed(1, func(gtx C) D {
											return style.Card{
												Content: []layout.Widget{
													func(gtx C) D {
														return material.H6(p.Th.Primary(), "Utilities").Layout(gtx)
													},
													func(gtx C) D {
														balance := 0
														if service, ok := p.lease.Services["utilities"]; ok {
															balance = service.Balance()
														}
														return style.ServiceLabel(p.Th, "Balance", float64(balance)).Layout(gtx)
													},
													func(gtx C) D {
														return layout.Flex{
															Axis:      layout.Horizontal,
															Alignment: layout.Middle,
														}.Layout(
															gtx,
															layout.Flexed(1, func(gtx C) D {
																b := material.Button(p.Th.Success(), &p.PayUtility, "Pay")
																b.Inset = layout.UniformInset(unit.Dp(5))
																return b.Layout(gtx)
															}),
															layout.Rigid(func(gtx C) D {
																return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
															}),
															layout.Flexed(1, func(gtx C) D {
																b := material.Button(p.Th.Danger(), &p.BillUtility, "Bill")
																b.Inset = layout.UniformInset(unit.Dp(5))
																return b.Layout(gtx)
															}),
														)
													},
												},
											}.Layout(gtx, p.Th.Primary())
										}),
										layout.Rigid(func(gtx C) D {
											return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
										}),
										layout.Flexed(1, func(gtx C) D {
											return style.Card{
												Content: []layout.Widget{
													func(gtx C) D {
														return material.H6(p.Th.Primary(), "Rent").Layout(gtx)
													},
													func(gtx C) D {
														balance := 0
														if service, ok := p.lease.Services["rent"]; ok {
															balance = service.Balance()
														}
														return style.ServiceLabel(p.Th, "Balance", float64(balance)).Layout(gtx)
													},
													func(gtx C) D {
														return layout.Flex{
															Axis:      layout.Horizontal,
															Alignment: layout.Middle,
														}.Layout(
															gtx,
															layout.Flexed(1, func(gtx C) D {
																b := material.Button(p.Th.Success(), &p.PayRent, "Pay")
																b.Inset = layout.UniformInset(unit.Dp(5))
																return b.Layout(gtx)
															}),
															layout.Rigid(func(gtx C) D {
																return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
															}),
															layout.Flexed(1, func(gtx C) D {
																b := material.Button(p.Th.Danger(), &p.BillRent, "Bill")
																b.Inset = layout.UniformInset(unit.Dp(5))
																return b.Layout(gtx)
															}),
														)
													},
												},
											}.Layout(gtx, p.Th.Primary())
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									return D{Size: image.Point{Y: gtx.Px(unit.Dp(10))}}
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							return material.Label(p.Th.Primary(), unit.Dp(20), "Details").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return D{Size: image.Point{Y: gtx.Px(unit.Dp(10))}}
						}),
						layout.Rigid(func(gtx C) D {
							return p.Form.Layout(gtx, p.Th)
						}),
					)
				}),
			)
		})
	})
}
