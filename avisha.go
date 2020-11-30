package avisha

import (
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/jackmordaunt/avisha-fn/notify"
)

// Tenant is a unique entity that can Lease one or more Sites.
type Tenant struct {
	Id      uint
	Name    string
	Contact string
}

// Site is a unique lot of land with a dwelling that can be Leased by at most
// one Tenant at any given time.
type Site struct {
	Id       uint
	Number   string
	Dwelling Dwelling
}

// Dwelling is where a Tenant lives.
type Dwelling int

const (
	// Cabin is a small, self contained temporary home.
	Cabin Dwelling = iota
	// Flat is a medium size permanent home.
	Flat
	// House is a full size permanent home.
	House
)

func (d Dwelling) String() string {
	switch d {
	case Flat:
		return "Flat"
	case House:
		return "House"
	case Cabin:
		return "Cabin"
	default:
		return "Unknown"
	}
}

// Currency in encoded in AUD cents, where 100 == $1.
type Currency = uint

// Term describes the active duration of a Lease.
type Term struct {
	Start time.Time
	Days  int
}

// Overlaps returns whether both terms overlap.
// A term overlaps another if they share active days.
func (t Term) Overlaps(other Term) bool {
	end := t.Start.Add(time.Hour * time.Duration(24) * time.Duration(t.Days))
	return other.Start.After(t.Start) && other.Start.Before(end) || other.Overlaps(t)
}

// Lease describes the exclusive use of a Site by exactly one Tenant for the
// duration of the Term specified.
// Services consumed are tracked accordingly, typically involving Rent and
// Utilities.
type Lease struct {
	Id     uint
	Tenant uint
	Site   uint

	Term Term
	Rent Currency

	// Rent and Utility services.
	Services map[string]Service
}

// Service is a billable for a lease.
// Typically Rent and Utilities are a services.
type Service struct {
	Credits []Currency
	Debits  []Currency
}

// Balance calculates the Balance of the Service.
func (s Service) Balance() int {
	credits := 0
	for _, c := range s.Credits {
		credits += int(c)
	}
	debits := 0
	for _, d := range s.Debits {
		debits += int(d)
	}
	return credits - debits
}

// App implements use cases.
type App struct {
	*storm.DB
	notify.Notifier
}

// SendInvoice sends the utility service invoice for the given Lease.
func (app App) SendInvoice(t Tenant, s Site) error {
	// containsSite := func(ent storage.Entity) bool {
	// 	if lease, ok := ent.(Lease); ok {
	// 		return lease.Site == s.Id
	// 	}
	// 	return false
	// }

	// containsTenant := func(ent storage.Entity) bool {
	// 	if lease, ok := ent.(Lease); ok {
	// 		return lease.Tenant == t.Id
	// 	}
	// 	return false
	// }

	// if entity, ok := app.Query(containsSite, containsTenant); ok {
	// 	if lease, ok := entity.(Lease); ok {

	// 		// Note: Instead of any actual invoice rendering we will just render
	// 		// utility balance owed.
	// 		if utilities, ok := lease.Services["utility"]; ok {
	// 			var (
	// 				balance = utilities.Balance()
	// 				dollars = balance % 100
	// 				cents   = balance - (dollars * 100)
	// 			)

	// 			if balance < 0 {
	// 				invoice := fmt.Sprintf("you owe $%2d.%2d in utilities", dollars, cents)
	// 				if err := app.Notify(t.Contact, invoice); err != nil {
	// 					return fmt.Errorf("sending invoice: %w", err)
	// 				}
	// 			}
	// 		}

	// 	}
	// }
	return nil
}

// CreateLease creates a new lease.
func (app App) CreateLease(
	tenant uint,
	site uint,
	term Term,
	rent Currency,
) error {
	// containsSite := func(ent storage.Entity) bool {
	// 	if lease, ok := ent.(*Lease); ok {
	// 		return lease.Site == site
	// 	}
	// 	return false
	// }

	// matchesTerm := func(ent storage.Entity) bool {
	// 	if lease, ok := ent.(*Lease); ok {
	// 		return lease.Term == term
	// 	}
	// 	return false
	// }

	// tenantExists := func(ent storage.Entity) bool {
	// 	if t, ok := ent.(*Tenant); ok {
	// 		if t.Id == tenant {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// siteExists := func(ent storage.Entity) bool {
	// 	if s, ok := ent.(*Site); ok {
	// 		if s.Id == site {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// if _, ok := app.Query(tenantExists); !ok {
	// 	return fmt.Errorf("tenant %d does not exist", tenant)
	// }

	// if _, ok := app.Query(siteExists); !ok {
	// 	return fmt.Errorf("site %d does not exist", site)
	// }

	// if _, ok := app.Query(containsSite, matchesTerm); ok {
	// 	return fmt.Errorf("lease conflict: site already leased during this term")
	// }

	// lease := Lease{
	// 	Tenant: tenant,
	// 	Site:   site,
	// 	Term:   term,
	// 	Rent:   rent,
	// 	Services: map[string]Service{
	// 		"rent":    {},
	// 		"utility": {},
	// 	},
	// }

	// if err := app.Create(lease); err != nil {
	// 	return fmt.Errorf("saving lease: %w", err)
	// }

	return nil
}

// ChangeRent updates the weekly rent for a given lease.
func (app App) ChangeRent(
	tenantID uint,
	siteID uint,
	term Term,
	rent Currency,
) error {
	// var (
	// 	lease  *Lease
	// 	tenant *Tenant
	// 	site   *Site
	// )

	// app.Get(0, func(ent storage.Entity) bool {
	// 	if t, ok := ent.(*Tenant); ok && t.Id == tenantID {
	// 		tenant = t
	// 		return true
	// 	}
	// 	return false
	// })

	// app.Get(0, func(ent storage.Entity) bool {
	// 	if s, ok := ent.(*Site); ok && s.Id == siteID {
	// 		site = s
	// 		return true
	// 	}
	// 	return false
	// })

	// app.Get(0, func(ent storage.Entity) bool {
	// 	if l, ok := ent.(*Lease); ok && l.Term == term {
	// 		lease = l
	// 		return true
	// 	}
	// 	return false
	// })

	// // if !ok {
	// // 	return fmt.Errorf("no lease exists for tenant %s and site %s during %s", tenantID, site, term)
	// // }

	// lease.Rent = rent
	// if err := app.Update(*lease); err != nil {
	// 	return fmt.Errorf("changing rent: %w", err)
	// }
	return nil
}

// ListSite enters a new, unqiue, leaseable Site.
func (app App) ListSite(s Site) error {
	// s.Number = strings.TrimSpace(s.Number)

	// exists := func(ent storage.Entity) bool {
	// 	if site, ok := ent.(*Site); ok {
	// 		return site.Number == s.Number
	// 	}
	// 	return false
	// }

	// if _, ok := app.Query(exists); ok {
	// 	return fmt.Errorf("%s already exists", s.Number)
	// }

	// if err := app.Create(s); err != nil {
	// 	return fmt.Errorf("saving site: %w", err)
	// }

	return nil
}

// RegisterTenant enters a new, unique Tenant.
func (app App) RegisterTenant(t Tenant) error {
	// exists := func(ent storage.Entity) bool {
	// 	if tenant, ok := ent.(*Tenant); ok {
	// 		return tenant.Name == t.Name
	// 	}
	// 	return false
	// }

	// if len(t.Name) < 1 {
	// 	return fmt.Errorf("name required")
	// }

	// if _, ok := app.Query(exists); ok {
	// 	return fmt.Errorf("%s already exists", t.Name)
	// }

	// t.Id = app.NextID()

	// if err := app.Create(t); err != nil {
	// 	return fmt.Errorf("saving tenant: %w", err)
	// }

	return nil
}

// ID specifies the unique identifier.
func (t Tenant) ID() uint {
	return t.Id
}

// ID specifies the unique identifier.
func (s Site) ID() uint {
	return s.Id
}

// ID specifies the unique identifier.
func (l Lease) ID() uint {
	return l.Id
}

func (t Term) String() string {
	format := func(t time.Time) string {
		return fmt.Sprintf("%02d/%02d/%04d", t.Day(), t.Month(), t.Year())
	}
	return fmt.Sprintf("%s - %s", format(t.Start), format(t.End()))
}

func (t Term) End() time.Time {
	return t.Start.Add(time.Hour * 24 * time.Duration(t.Days))
}

// LeaseComparitor can be used for comparison between lease entities.
type LeaseComparitor struct {
	Tenant uint
	Site   uint
	Term   Term
}

func (l Lease) Cmp() LeaseComparitor {
	return LeaseComparitor{
		Tenant: l.Tenant,
		Site:   l.Site,
		Term:   l.Term,
	}
}
