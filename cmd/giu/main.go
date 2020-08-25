package main

import (
	"fmt"
	"time"

	"github.com/AllenDang/giu"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"
)

func main() {
	app := avisha.App{
		Storage: storage.FileStorage("target/db.json").
			With(&avisha.Tenant{}).
			With(&avisha.Site{}).
			With(&avisha.Lease{}).
			Format(true).
			MustLoad(),
		Notifier: &notify.Console{},
	}

	state := &State{
		App: app,
	}

	giu.NewMasterWindow(
		"Avisha",
		1024,
		768,
		0,
		nil,
	).Main(func() {
		giu.SingleWindow("Avisha", giu.Layout{
			state,
		})
	})
}

// State contains the state of the gui.
type State struct {
	// App is the application.
	avisha.App
	// Page is the currently presented page.
	Page giu.Widget
}

// Build the gui.
func (s *State) Build() {
	if s.Page == nil {
		s.Page = s.DashBoardPage()
	}
	giu.SplitLayout("Main", giu.DirectionHorizontal, true, 200,
		// Nav sidebar.
		giu.Layout{
			giu.ButtonV("DashBoard", -1, 40, func() {
				s.Page = s.DashBoardPage()
			}),
			giu.ButtonV("Lease", -1, 40, func() {
				s.Page = s.LeasePage()
			}),
			giu.ButtonV("Tenant", -1, 40, func() {
				s.Page = s.TenantPage()
			}),
			giu.ButtonV("Site", -1, 40, func() {
				s.Page = s.SitePage()
			}),
		},
		s.Page,
	).Build()
}

// DashBoardPage renders the Dashboard.
func (s *State) DashBoardPage() giu.Widget {
	var labels giu.Layout
	s.List(func(ent storage.Entity) bool {
		l, ok := ent.(*avisha.Lease)
		if ok {
			labels = append(labels, giu.Label(l.ID()))
		}
		return ok
	})
	return labels
}

// LeasePage renders Lease controls.
func (s *State) LeasePage() giu.Widget {
	return &LeaseForm{
		App: s.App,
		OnSubmit: func(f *LeaseForm) {
			err := s.App.CreateLease(f.Tenant, f.Site, avisha.Term{
				Start: f.Start,
				Days:  int(f.Days),
			}, uint(f.Rent))
			if err != nil {
				fmt.Printf("CreateLease: %v\n", err)
			}
		},
	}
}

// TenantPage renders Tenant controls.
func (s *State) TenantPage() giu.Widget {
	return &TenantForm{
		OnSubmit: func(f *TenantForm) {
			err := s.App.RegisterTenant(avisha.Tenant{
				Name:    f.Name,
				Contact: f.Contact,
			})
			if err != nil {
				fmt.Printf("RegisterTenant: %v\n", err)
			}
		},
	}
}

// SitePage renders site controls.
func (s *State) SitePage() giu.Widget {
	return &SiteForm{
		OnSubmit: func(f *SiteForm) {
			err := s.App.ListSite(avisha.Site{
				Number:   f.Number,
				Dwelling: f.Dwelling,
			})
			if err != nil {
				fmt.Printf("ListSite: %v\n", err)
			}
		},
	}
}

// LeaseForm for creating a Lease.
type LeaseForm struct {
	avisha.App

	Tenant string
	Site   string
	Start  time.Time
	Days   int32
	Rent   int32

	OnSubmit func(*LeaseForm)

	tenant *Combo
	site   *Combo
}

// Build LeaseForm
func (form *LeaseForm) Build() {
	if form.tenant == nil {
		form.tenant = &Combo{
			Label:   "Tenant Name",
			Preview: "Pick a Tenant",
			OnChange: func(t string) {
				form.Tenant = t
			},
		}
		form.App.List(func(ent storage.Entity) bool {
			t, ok := ent.(*avisha.Tenant)
			if ok {
				form.tenant.Add(t.Name, t.Name)
			}
			return ok
		})
	}

	if form.site == nil {
		form.site = &Combo{
			Label:   "Site Number",
			Preview: "Pick a Site",
			OnChange: func(s string) {
				form.Site = s
			},
		}
		form.App.List(func(ent storage.Entity) bool {
			s, ok := ent.(*avisha.Site)
			if ok {
				form.site.Add(fmt.Sprintf("%s - %s", s.Number, s.Dwelling), s.Number)
			}
			return ok
		})
	}

	giu.Layout{
		form.tenant,
		form.site,
		giu.DatePicker("Start", &form.Start, -1, func() {}),
		giu.InputInt("Days", -1, &form.Days),
		giu.InputInt("Rent", -1, &form.Rent),
		giu.Button("Submit", func() {
			if form.OnSubmit != nil {
				form.OnSubmit(form)
			}
			form.tenant = nil
			form.site = nil
			form.Tenant = ""
			form.Site = ""
			form.Start = time.Now()
			form.Days = 0
			form.Rent = 0
		}),
	}.Build()
}

// Combo renders a combobox.
type Combo struct {
	Label string
	// Preview to display when nothing has been selected.
	Preview string
	// Items holds the formatted string to display.
	Items []string
	// Values holds the value to pass to OnChange.
	Values   []string
	OnChange func(selected string)

	selected    int32
	hasSelected bool
}

// Add an item:value pair to the combo.
func (c *Combo) Add(item, value string) {
	c.Items = append(c.Items, item)
	c.Values = append(c.Values, value)
}

// Build Combo
func (c *Combo) Build() {
	preview := func() string {
		if !c.hasSelected {
			return c.Preview
		}
		return c.item()
	}
	giu.Combo(c.Label, preview(), c.Items, &c.selected, -1, 0, func() {
		c.hasSelected = true
		if c.OnChange != nil {
			c.OnChange(c.value())
		}
	}).Build()
}

func (c *Combo) value() string {
	if int(c.selected) < len(c.Values) {
		return c.Values[c.selected]
	}
	return ""
}

func (c *Combo) item() string {
	if int(c.selected) < len(c.Items) {
		return c.Items[c.selected]
	}
	return ""
}

// Date renders a date picker.
type Date struct {
	Label    string
	OnChange func(*time.Time)

	value *time.Time
}

// Build Date
func (d *Date) Build() {
	if d.value == nil {
		now := time.Now()
		d.value = &now
	}
	giu.DatePicker(d.Label, d.value, -1, func() {
		if d.OnChange != nil {
			d.OnChange(d.value)
		}
	}).Build()
}

// SiteForm creates a Site.
type SiteForm struct {
	Number   string
	Dwelling avisha.Dwelling

	OnSubmit func(*SiteForm)

	dwelling *Combo
}

// Build SiteForm
func (form *SiteForm) Build() {
	if form.dwelling == nil {
		form.dwelling = &Combo{
			Label:   "Dwelling",
			Preview: "Pick a Dwelling",
			Items:   []string{"House", "Cabin", "Flat"},
			OnChange: func(dwelling string) {
				switch dwelling {
				case "House":
					form.Dwelling = avisha.House
				case "Cabin":
					form.Dwelling = avisha.Cabin
				case "Flat":
					form.Dwelling = avisha.Flat
				default:
					form.Dwelling = avisha.House
				}
			},
		}
	}
	giu.Layout{
		giu.InputText("Number", -1, &form.Number),
		form.dwelling,
		giu.Button("Submit", func() {
			if form.OnSubmit != nil {
				form.OnSubmit(form)
				form.Number = ""
				form.Dwelling = avisha.House
			}
		}),
	}.Build()
}

// TenantForm inputs Tenant data.
type TenantForm struct {
	Name    string
	Contact string

	OnSubmit func(*TenantForm)
}

// Build TenantForm.
func (form *TenantForm) Build() {
	giu.Layout{
		giu.InputText("Name", -1, &form.Name),
		giu.InputText("Contact", -1, &form.Contact),
		giu.Button("Submit", func() {
			if form.OnSubmit != nil {
				form.OnSubmit(form)
				form.Name = ""
				form.Contact = ""
			}
		}),
	}.Build()
}
