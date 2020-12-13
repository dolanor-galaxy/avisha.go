package avisha

import (
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/jackmordaunt/avisha-fn/notify"
)

type ID = int

// Tenant is a unique entity that can Lease one or more Sites.
type Tenant struct {
	ID      ID     `storm:"id,increment"`
	Name    string `storm:"unique"`
	Contact string
}

// Site is a unique lot of land with a dwelling that can be Leased by at most
// one Tenant at any given time.
type Site struct {
	ID       ID     `storm:"id,increment"`
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
	ID     ID `storm:"id,increment" `
	Tenant int
	Site   int

	Term Term
	Rent Currency

	// Services is a map of named services like rent and utilities.
	Services map[string]Service
}

// Service is a billable for a lease.
type Service struct {
	Ledger Ledger
}

func (s Service) Balance() int {
	return s.Ledger.Balance()
}

// Ledger maintains a balance of currency.
type Ledger struct {
	Credits []Currency
	Debits  []Currency
}

// Credit record a credit of currency.
func (l *Ledger) Credit(amount Currency) {
	l.Credits = append(l.Credits, amount)
}

// Debit records a debit of currency.
func (l *Ledger) Debit(amount Currency) {
	l.Debits = append(l.Debits, amount)
}

// Balance calculates the Balance of the Service.
func (l Ledger) Balance() int {
	credits := 0
	for _, c := range l.Credits {
		credits += int(c)
	}
	debits := 0
	for _, d := range l.Debits {
		debits += int(d)
	}
	return credits - debits
}

// Invoice is a document requesting payment for a service.
type Invoice struct {
	ID    ID `storm:"id,increment"`
	Lease int
	// Bill is the amount of currency due.
	Bill Currency
	// Important dates.
	Issued time.Time
	Due    time.Time
	Paid   time.Time
}

// UtilityInvoice is a document requesting payment for utility consumption.
type UtilityInvoice struct {
	Invoice `storm:"inline"`
	// UnitCost is the cost per unit of power.
	UnitCost Currency
	// UnitsConsumed is the amount of units to charge for.
	UnitsConsumed int
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

// PayService records a payment for some service on a lease.
func (app App) PayService(leaseID int, service string, amount uint) error {
	var l Lease
	if err := app.One("ID", leaseID, &l); err != nil {
		return fmt.Errorf("finding lease: %w", err)
	}
	if l.Services == nil {
		l.Services = make(map[string]Service)
	}
	s := l.Services[service]
	s.Ledger.Credit(amount)
	l.Services[service] = s
	if err := app.markInvoices(leaseID, s); err != nil {
		return fmt.Errorf("marking invoices: %w", err)
	}
	return app.Update(&l)
}

// BillService records a debt for some service on a lease.
func (app App) BillService(leaseID int, service string, amount uint) error {
	var l Lease
	if err := app.One("ID", leaseID, &l); err != nil {
		return fmt.Errorf("finding lease: %w", err)
	}
	if l.Services == nil {
		l.Services = make(map[string]Service)
	}
	s := l.Services[service]
	s.Ledger.Debit(amount)
	l.Services[service] = s
	if err := app.markInvoices(leaseID, s); err != nil {
		return fmt.Errorf("marking invoices: %w", err)
	}
	return app.Update(&l)
}

// markInvoices marks invoices for a given service as paid, starting from oldest
// first.
func (app App) markInvoices(leaseID int, service Service) error {
	var (
		total    int
		invoices []*UtilityInvoice
	)
	if err := app.Select(q.Eq("Lease", leaseID)).OrderBy("ID").Find(&invoices); err != nil {
		return fmt.Errorf("loading invoices: %w", err)
	}
	for _, credit := range service.Ledger.Credits {
		total += int(credit)
	}
	// Pay all the invoices we can, marking them paid as of now if they weren't
	// marked already.
	//
	// @Fix this is kinda hacky. Think about how payments and invoices should
	// interact.
	for _, inv := range invoices {
		if total < int(inv.Bill) {
			break
		}
		if inv.Paid == (time.Time{}) {
			inv.Paid = time.Now()
		}
		if err := app.Update(inv); err != nil {
			return fmt.Errorf("update: %w", err)
		}
		total -= int(inv.Bill)
	}
	return nil
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
