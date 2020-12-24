package views

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/nav"
	"github.com/jackmordaunt/avisha.go/cmd/gui/util"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
	"github.com/jackmordaunt/avisha.go/currency"
	"github.com/skratchdot/open-golang/open"
)

// LeasePage contains actions for interacting with a lease including data entry
// and service payments.
type LeasePage struct {
	nav.Route
	App *avisha.App
	Th  *style.Theme

	lease avisha.Lease

	Form                 LeaseForm
	Dialog               style.Dialog
	UtilitiesInvoiceForm UtilitiesInvoiceForm

	PayUtility  widget.Clickable
	BillUtility widget.Clickable
	PayRent     widget.Clickable
	BillRent    widget.Clickable

	modal         layout.Widget
	invoiceStates States
	invoiceList   layout.List
	scroll        layout.List
	dummy         widget.Editor
}

func (page *LeasePage) Title() string {
	return "Lease"
}

func (p *LeasePage) Receive(data interface{}) {
	p.Form.App = p.App
	if lease, ok := data.(*avisha.Lease); ok && lease != nil {
		p.lease = *lease
	} else {
		p.lease = avisha.Lease{}
		defer p.Form.Clear()
	}
	p.Form.Load(p.lease)
	p.UtilitiesInvoiceForm.Clear()
}

func (p *LeasePage) Context() (list []layout.Widget) {
	if p.lease.ID != 0 {
		list = append(list, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(
				gtx,
				func(gtx C) D {
					label := material.Label(
						p.Th.Dark(),
						unit.Dp(24),
						fmt.Sprintf("%s-%s", p.Form.Site.Text(), p.Form.Tenant.Text()))
					label.Alignment = text.Middle
					label.Color = p.Th.ContrastFg
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
	if p.Form.SubmitBtn.Clicked() {
		if lease, ok := p.Form.Submit(); ok {
			if err := func() error {
				if create := p.lease.ID == 0; create {
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
				p.Unfocus()
				if err := p.App.One("ID", p.lease.ID, &p.lease); err != nil {
					log.Printf("error: loading lease: %d: %v", p.lease.ID, err)
				}
			}
		}
	}
	if p.Form.CancelBtn.Clicked() {
		p.Unfocus()
		p.Form.Load(p.lease)
		if p.lease.ID == 0 {
			p.Form.Clear()
			p.Route.Back()
		}
	}
	if p.PayUtility.Clicked() {
		p.Dialog.Context = "pay-utilities"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Pay Utilities", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Dark(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	if p.BillUtility.Clicked() {
		p.Dialog.Context = "bill-utilities"
		var (
			prevReading = 0
			invoices    []*avisha.UtilityInvoice
		)
		if err := p.App.Select(q.Eq("Lease", p.lease.ID)).OrderBy("ID", "Paid").Reverse().Find(&invoices); err != nil {
			if err != storm.ErrNotFound {
				log.Printf("loading invoices: %v", err)
			}
		}
		if len(invoices) > 0 {
			prevReading = invoices[len(invoices)-1].Reading
		}
		p.UtilitiesInvoiceForm.Load(avisha.UtilityInvoice{}, prevReading)
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Bill Utilities", func(gtx C) D {
				return p.UtilitiesInvoiceForm.Layout(gtx, p.Th)
			})
		}
	}
	if p.PayRent.Clicked() {
		p.Dialog.Context = "pay-rent"
		p.modal = func(gtx C) D {
			return style.ModalDialog(gtx, p.Th, unit.Dp(700), "Pay Rent", func(gtx C) D {
				p.Dialog.Input.Prefix = func(gtx C) D {
					return material.Label(p.Th.Dark(), p.Th.TextSize, "$").Layout(gtx)
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
					return material.Label(p.Th.Dark(), p.Th.TextSize, "$").Layout(gtx)
				}
				return p.Dialog.Layout(gtx, p.Th.Primary(), "Amount")
			})
		}
	}
	for range p.Dialog.Input.Events() {
		_, err := util.ParseCurrency(p.Dialog.Input.Text())
		if err != nil {
			p.Dialog.Input.SetError(err.Error())
		} else {
			p.Dialog.Input.ClearError()
		}
	}
	if p.Dialog.Ok.Clicked() {
		if n, err := util.ParseCurrency(p.Dialog.Input.Text()); err != nil {
			p.Dialog.Input.SetError(err.Error())
		} else {
			p.Dialog.Input.Clear()
			// @Improvment Can we avoid stringly typed api?
			parts := strings.Split(p.Dialog.Context, "-")
			mode, service := parts[0], parts[1]
			switch mode {
			case "pay":
				if err := p.App.PayService(p.lease.ID, service, n); err != nil {
					log.Printf("paying service: %v", err)
				}
			case "bill":
				if err := p.App.BillService(p.lease.ID, service, n); err != nil {
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
	if p.UtilitiesInvoiceForm.SubmitBtn.Clicked() {
		if inv, ok := p.UtilitiesInvoiceForm.Submit(); ok {
			inv.Lease = p.lease.ID
			if err := p.App.Save(&inv); err != nil {
				log.Printf("saving invoice: %v", err)
			}
			if err := p.App.BillService(
				p.lease.ID,
				"utilities",
				inv.UnitCost*currency.Currency(inv.UnitsConsumed),
			); err != nil {
				log.Printf("billing service: %v", err)
			}
		}
		p.UtilitiesInvoiceForm.Clear()
		p.modal = nil
	}
	if p.UtilitiesInvoiceForm.CancelBtn.Clicked() {
		p.UtilitiesInvoiceForm.Clear()
		p.modal = nil
	}
	for _, state := range p.invoiceStates.List() {
		var (
			invoice = (*avisha.UtilityInvoice)(state.Data)
		)
		// @Todo Do io async to avoid blocking ui.
		if state.Item.Clicked() {
			if err := func() error {
				doc := util.UtilityInvoiceDocument{
					Invoice: *invoice,
				}
				buffer, err := doc.Render()
				if err != nil {
					return fmt.Errorf("rendering invoice document: %w", err)
				}
				dir, err := app.DataDir()
				if err != nil {
					return fmt.Errorf("locating data directory: %w", err)
				}
				dir = filepath.Join(dir, "invoices")
				if err := os.MkdirAll(dir, 0777); err != nil && !os.IsNotExist(err) {
					return fmt.Errorf("preparing directory: %w", err)
				}
				path := filepath.Join(dir, fmt.Sprintf("%d.html", invoice.ID))
				if err := ioutil.WriteFile(
					path,
					buffer.Bytes(),
					0777,
				); err != nil {
					return fmt.Errorf("writing invoice to disk: %w", err)
				}
				if err := open.Run(path); err != nil {
					return fmt.Errorf("opening invoice: %w", err)
				}
				return nil
			}(); err != nil {
				log.Printf("%v", err)
			}
		}
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
	p.scroll.Axis = layout.Vertical
	p.scroll.ScrollToEnd = false
	p.Update(gtx)
	var (
		cs    = &gtx.Constraints
		inset = layout.UniformInset(unit.Dp(10))
		axis  = layout.Vertical
	)
	if breakpoint := gtx.Px(unit.Dp(800)); cs.Max.X > breakpoint {
		if p.lease.ID > 0 {
			axis = layout.Horizontal
		} else {
			// Cap the form's max width on expanded view.
			cs.Max.X = breakpoint
		}
	}
	func() {
		// @Improvement proper strategy for handling "unfocus".
		// Currenty strategy is to use a dummy editor to take the focus away.
		// Note: could use the dummy for capturing commands.
		th := material.Theme{
			Shaper:   p.Th.Shaper,
			TextSize: unit.Dp(0),
		}
		material.Editor(&th, &p.dummy, "").Layout(gtx)
	}()
	return p.scroll.Layout(gtx, 1, func(gtx C, ii int) D {
		return layout.Flex{
			Axis:      axis,
			Alignment: layout.Start,
		}.Layout(
			gtx,
			util.FlexStrategy(1, layout.Horizontal, axis, func(gtx C) D {
				if p.lease.ID == 0 {
					return D{}
				}
				return inset.Layout(gtx, func(gtx C) D {
					return p.LayoutServices(gtx)
				})
			}),
			util.FlexStrategy(1, layout.Horizontal, axis, func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					return p.LayoutDetails(gtx)
				})
			}),
			util.FlexStrategy(1, layout.Horizontal, axis, func(gtx C) D {
				if p.lease.ID == 0 {
					return D{}
				}
				return inset.Layout(gtx, func(gtx C) D {
					return p.LayoutInvoiceList(gtx)
				})
			}),
		)
	})
}

// LayoutDetails lays the details form.
func (p *LeasePage) LayoutDetails(gtx C) D {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return material.Label(p.Th.Dark(), unit.Dp(20), "Details").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{Y: gtx.Px(unit.Dp(10))}}
		}),
		layout.Rigid(func(gtx C) D {
			return p.Form.Layout(gtx, p.Th)
		}),
	)
}

// LayoutServices lays the service widgets in flex.
func (p *LeasePage) LayoutServices(gtx C) D {
	var (
		cs   = &gtx.Constraints
		axis = layout.Horizontal
	)
	if breakpoint := gtx.Px(unit.Dp(350)); cs.Max.X < breakpoint {
		axis = layout.Vertical
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return material.Label(p.Th.Dark(), unit.Dp(20), "Services").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10)), Y: gtx.Px(unit.Dp(10))}}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis:      axis,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Start,
			}.Layout(
				gtx,
				util.FlexStrategy(1, layout.Horizontal, axis, func(gtx C) D {
					return style.Card{
						Content: []layout.Widget{
							func(gtx C) D {
								return material.H6(p.Th.Dark(), "Utilities").Layout(gtx)
							},
							func(gtx C) D {
								var balance currency.Currency
								if service, ok := p.lease.Services["utilities"]; ok {
									balance = service.Balance()
								}
								return style.ServiceLabel(p.Th, "Balance", balance).Layout(gtx)
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
					}.Layout(gtx, p.Th.Dark())
				}),
				layout.Rigid(func(gtx C) D {
					return D{Size: image.Point{X: gtx.Px(unit.Dp(10)), Y: gtx.Px(unit.Dp(10))}}
				}),
				util.FlexStrategy(1, layout.Horizontal, axis, func(gtx C) D {
					return style.Card{
						Content: []layout.Widget{
							func(gtx C) D {
								return material.H6(p.Th.Dark(), "Rent").Layout(gtx)
							},
							func(gtx C) D {
								var balance currency.Currency
								if service, ok := p.lease.Services["rent"]; ok {
									balance = service.Balance()
								}
								return style.ServiceLabel(p.Th, "Balance", balance).Layout(gtx)
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
					}.Layout(gtx, p.Th.Dark())
				}),
			)
		}),
	)
}

// LayoutInvoiceList renders a list of invoices issued for the lease.
func (p *LeasePage) LayoutInvoiceList(gtx C) D {
	p.invoiceList.Axis = layout.Vertical
	p.invoiceList.ScrollToEnd = false
	p.invoiceStates.Begin()
	var (
		invoices []*avisha.UtilityInvoice
	)
	if err := p.App.Select(q.Eq("Lease", p.lease.ID)).OrderBy("ID", "Paid").Reverse().Find(&invoices); err != nil {
		if err != storm.ErrNotFound {
			log.Printf("loading invoices: %v", err)
		}
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return material.Label(p.Th.Dark(), unit.Dp(20), "Invoices").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10)), Y: gtx.Px(unit.Dp(10))}}
		}),
		layout.Rigid(func(gtx C) D {
			return p.invoiceList.Layout(gtx, len(invoices), func(gtx C, ii int) D {
				var (
					invoice = invoices[ii]
					state   = p.invoiceStates.Next(unsafe.Pointer(invoice))
					active  = false
				)
				return style.ListItem(
					gtx,
					p.Th.Dark(),
					&state.Item,
					&state.Hover,
					active,
					func(gtx C) D {
						return layout.Flex{
							Axis: layout.Horizontal,
						}.Layout(
							gtx,
							layout.Flexed(3, func(gtx C) D {
								return material.Label(
									p.Th.Dark(),
									unit.Dp(14),
									fmt.Sprintf(
										"#%d %s (%d %s %d)",
										invoice.ID,
										invoice.Bill,
										invoice.Issued.Day(),
										invoice.Issued.Month(),
										invoice.Issued.Year()),
								).Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								var (
									badge = "PAID"
									c     = p.Th.Success().Fg
								)
								if invoice.Paid == (time.Time{}) {
									badge = "NOT PAID"
									c = p.Th.Danger().Fg
								}
								lb := material.Label(
									p.Th.Dark(),
									unit.Dp(14),
									badge,
								)
								lb.Color = c
								return lb.Layout(gtx)
							}),
						)
					},
				)
			})
		}),
	)
}

func (p *LeasePage) Unfocus() {
	p.dummy.Focus()
}

// TenantValuer maps tenant name to ID.
//
// @Todo How to handle several tenants with the same name?
type TenantValuer struct {
	ID  *int
	App *avisha.App
}

func (v TenantValuer) To() (string, error) {
	var t avisha.Tenant
	if err := v.App.One("ID", *v.ID, &t); err != nil {
		return "", err
	}
	return t.Name, nil
}

func (v TenantValuer) From(text string) error {
	var t avisha.Tenant
	if err := v.App.One("Name", text, &t); err != nil {
		return err
	}
	*v.ID = t.ID
	return nil
}

func (v TenantValuer) Clear() {
	*v.ID = -1
}

// SiteValuer maps site number to ID.
//
// @Todo How to handle several sites with the same number?
type SiteValuer struct {
	ID  *int
	App *avisha.App
}

func (v SiteValuer) To() (string, error) {
	var s avisha.Site
	if err := v.App.One("ID", *v.ID, &s); err != nil {
		return "", err
	}
	return s.Number, nil
}

func (v SiteValuer) From(text string) error {
	var s avisha.Site
	if err := v.App.One("Number", text, &s); err != nil {
		return err
	}
	*v.ID = s.ID
	return nil
}

func (v SiteValuer) Clear() {
	*v.ID = -1
}
