package avisha

import (
	"fmt"
	"time"

	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"
)

// Tenant is a unique entity that can Lease one or more Sites.
type Tenant struct {
	Name    string
	Contact string
}

// Site is a unique lot of land with a dwelling that can be Leased by at most
// one Tenant at any given time.
type Site struct {
	Number   string
	Dwelling Dwelling
}

// Dwelling is where a Tenant lives.
type Dwelling = int

const (
	// Cabin is a small, self contained temporary home.
	Cabin Dwelling = iota
	// Flat is a medium size permanent home.
	Flat
	// House is a full size permanent home.
	House
)

// Currency in encoded in AUD cents, where 100 == $1.
type Currency = uint

// Term describes the active duration of a Lease.
type Term struct {
	Start time.Time
	Days  int
}

// Lease describes the exclusive use of a Site by exactly one Tenant for the
// duration of the Term specified.
// Services consumed are tracked accordingly, typically involving Rent and
// Utilities.
type Lease struct {
	Tenant Tenant
	Site   Site
	Term   Term
	Rent   Currency

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
	storage.Storer
	notify.Notifier
}

// SendInvoice sends the utility service invoice for the given Lease.
func (app App) SendInvoice(t Tenant, s Site) error {
	containsSite := storage.PredicateFunc(func(ent interface{}) bool {
		if lease, ok := ent.(Lease); ok {
			return lease.Site.Number == s.Number
		}
		return false
	})

	containsTenant := storage.PredicateFunc(func(ent interface{}) bool {
		if lease, ok := ent.(Lease); ok {
			return lease.Tenant.Name == t.Name
		}
		return false
	})

	if entity, ok := app.Query([]storage.Predicate{containsSite, containsTenant}); ok {
		if lease, ok := entity.(Lease); ok {

			// Note: Instead of any actual invoice rendering we will just render
			// utility balance owed.
			if utilities, ok := lease.Services["utility"]; ok {
				var (
					balance = utilities.Balance()
					dollars = balance % 100
					cents   = balance - (dollars * 100)
				)

				if balance < 0 {
					invoice := fmt.Sprintf("you owe $%2d.%2d in utilities", dollars, cents)
					if err := app.Notify(t.Contact, invoice); err != nil {
						return fmt.Errorf("sending invoice: %w", err)
					}
				}
			}

		}
	}
	return nil
}

// CreateLease creates a new lease.
func (app App) CreateLease(
	t Tenant,
	s Site,
	term Term,
	rent Currency,
) error {
	containsSite := storage.PredicateFunc(func(ent interface{}) bool {
		if lease, ok := ent.(Lease); ok {
			return lease.Site.Number == s.Number
		}
		return false
	})

	matchesTerm := storage.PredicateFunc(func(ent interface{}) bool {
		if lease, ok := ent.(Lease); ok {
			return lease.Term == term
		}
		return false
	})

	if _, ok := app.Query([]storage.Predicate{containsSite, matchesTerm}); ok {
		return fmt.Errorf("lease conflict: site already leased during this term")
	}

	lease := Lease{
		Tenant: t,
		Site:   s,
		Term:   term,
		Rent:   rent,
		Services: map[string]Service{
			"rent":    {},
			"utility": {},
		},
	}

	if err := app.Save(lease); err != nil {
		return fmt.Errorf("saving lease: %w", err)
	}

	return nil
}

// ListSite enters a new, unqiue, leaseable Site.
func (app App) ListSite(s Site) error {
	exists := storage.PredicateFunc(func(ent interface{}) bool {
		if site, ok := ent.(Site); ok {
			return site.Number == s.Number
		}
		return false
	})

	if _, ok := app.Query([]storage.Predicate{exists}); ok {
		return fmt.Errorf("%s already exists", s.Number)
	}

	if err := app.Save(s); err != nil {
		return fmt.Errorf("saving site: %w", err)
	}

	return nil
}

// RegisterTenant enters a new, unique Tenant.
func (app App) RegisterTenant(t Tenant) error {
	exists := storage.PredicateFunc(func(ent interface{}) bool {
		if tenant, ok := ent.(Tenant); ok {
			return tenant.Name == t.Name
		}
		return false
	})

	if len(t.Name) < 1 {
		return fmt.Errorf("name required")
	}

	if _, ok := app.Query([]storage.Predicate{exists}); ok {
		return fmt.Errorf("%s already exists", t.Name)
	}

	if err := app.Save(t); err != nil {
		return fmt.Errorf("saving tenant: %w", err)
	}

	return nil
}
