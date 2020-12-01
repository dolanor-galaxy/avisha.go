package avisha

import (
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/jackmordaunt/avisha-fn/notify"
)

// Tenant is a unique entity that can Lease one or more Sites.
type Tenant struct {
	ID      int    `storm:"id,increment"`
	Name    string `storm:"unique"`
	Contact string
}

// Site is a unique lot of land with a dwelling that can be Leased by at most
// one Tenant at any given time.
type Site struct {
	ID       int    `storm:"id,increment"`
	Number   string `storm:"unique"`
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
	ID     int `storm:"id,increment" `
	Tenant int
	Site   int

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

// CreateLease creates a new lease.
func (app App) CreateLease(l *Lease) error {
	if l.Tenant == 0 {
		return fmt.Errorf("lease must have a valid tenant")
	}
	if l.Site == 0 {
		return fmt.Errorf("lease must have a valid site")
	}
	return app.Save(l)
}

// ListSite enters a new, unqiue, leaseable Site.
func (app App) ListSite(s *Site) error {
	if len(s.Number) < 1 {
		return fmt.Errorf("number required")
	}
	return app.Save(s)
}

// RegisterTenant enters a new, unique Tenant.
func (app App) RegisterTenant(t *Tenant) error {
	if len(t.Name) < 1 {
		return fmt.Errorf("name required")
	}
	return app.Save(t)
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
	Tenant int
	Site   int
	Term   Term
}

func (l Lease) Cmp() LeaseComparitor {
	return LeaseComparitor{
		Tenant: l.Tenant,
		Site:   l.Site,
		Term:   l.Term,
	}
}
