package avisha

import (
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/jackmordaunt/avisha.go/currency"
	"github.com/jackmordaunt/avisha.go/notify"
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

// Term describes the active duration of a Lease.
type Term struct {
	Start    time.Time
	Duration time.Duration
}

// Overlaps returns whether both terms overlap.
// A term overlaps another if they share active days.
// @Todo test this.
func (t Term) Overlaps(other Term) bool {
	end := t.Start.Add(time.Duration(t.Duration))
	return other.Start.After(t.Start) && other.Start.Before(end) || other.Overlaps(t)
}

// Lease describes the exclusive use of a Site by exactly one Tenant for the
// duration of the Term specified.
// Services consumed are tracked accordingly, typically involving Rent and
// Utilities.
type Lease struct {
	ID     ID `storm:"id,increment"`
	Tenant int
	Site   int

	Term Term
	Rent currency.Currency

	// Services is a map of named services like rent and utilities.
	Services map[string]Service
}

// Service is a billable for a lease.
type Service struct {
	Ledger Ledger
}

func (s Service) Balance() currency.Currency {
	return s.Ledger.Balance()
}

// Ledger maintains a balance of currency.currency.
type Ledger struct {
	Credits []currency.Currency
	Debits  []currency.Currency
}

// Credit record a credit of currency.currency.
func (l *Ledger) Credit(amount currency.Currency) {
	l.Credits = append(l.Credits, amount)
}

// Debit records a debit of currency.currency.
func (l *Ledger) Debit(amount currency.Currency) {
	l.Debits = append(l.Debits, amount)
}

// Balance calculates the Balance of the Service.
func (l Ledger) Balance() currency.Currency {
	credits := currency.Currency(0)
	for _, c := range l.Credits {
		credits += c
	}
	debits := currency.Currency(0)
	for _, d := range l.Debits {
		debits += d
	}
	return credits - debits
}

// Invoice is a document requesting payment for a service.
type Invoice struct {
	ID    ID `storm:"id,increment"`
	Lease int
	// Bill is the amount of currency.currency due.
	Bill currency.Currency
	// Important dates.
	Issued time.Time
	Due    time.Time
	Paid   time.Time
	Period Term
}

// UtilityInvoice is a document requesting payment for utility consumption.
type UtilityInvoice struct {
	Invoice `storm:"inline"`
	// UnitCost is the cost per unit of power.
	UnitCost currency.Currency
	// UnitsConsumed is the amount of units to charge for.
	UnitsConsumed int
	// Reading is the units read off the meter.
	Reading int
	// GST records the GST used at the time the invoice was generated.
	GST float64
	// Charges contains all the constituent parts of the total bill.
	Charges struct {
		// LateFee for when payment is late.
		LateFee currency.Currency
		// LineCharge fee.
		LineCharge currency.Currency
		// GST calculated based on percentage.
		GST currency.Currency
		// Activity charge is "units-consumed * unit-cost".
		Activity currency.Currency
	}
}

// Settings are global settings that don't pertain to any specific entity.
type Settings struct {
	Landlord Landlord
	Bank     Bank
	Defaults Defaults
}

// Landlord details.
type Landlord struct {
	Name  string
	Email string
	Phone string
}

// Banks details to make invoices payable to.
type Bank struct {
	Name    string
	Account string
}

type Defaults struct {
	UnitCost   currency.Currency
	RentCycle  time.Duration
	InvoiceNet time.Duration
	GST        float64
}

// Default to sane values.
func (d *Defaults) Default() {
	d.UnitCost = 1
	d.RentCycle = time.Hour * 24 * 14
	d.InvoiceNet = time.Hour * 24 * 14
}

// App implements use cases.
type App struct {
	*storm.DB
	notify.Notifier
}

// LoadSettings loads global settings.
func (app App) LoadSettings() (s Settings, err error) {
	if err := app.Get("settings", "global", &s); err != nil && err != storm.ErrNotFound {
		return s, err
	}
	if s.Defaults == (Defaults{}) {
		s.Defaults.Default()
	}
	return s, nil
}

// SaveSettings saves global settings.
func (app App) SaveSettings(s Settings) (err error) {
	if s.Defaults == (Defaults{}) {
		s.Defaults.Default()
	}
	return app.Set("settings", "global", &s)
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
func (app App) PayService(leaseID int, service string, amount currency.Currency) error {
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
func (app App) BillService(leaseID int, service string, amount currency.Currency) error {
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
	// @Note this is valid for utilties, but rent needs to be payable out-of-order.
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
	return fmt.Sprintf("%s to %s", format(t.Start), format(t.End()))
}

func (t Term) End() time.Time {
	return t.Start.Add(t.Duration)
}
