package avisha

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/go-test/deep"
	"github.com/jackmordaunt/avisha.go/currency"
)

type mutexCounter struct {
	sync.Mutex
	Count int
}

func (m *mutexCounter) Next() int {
	m.Lock()
	defer m.Unlock()
	m.Count++
	return m.Count
}

var counter = mutexCounter{}

func OpenDB(t *testing.T) (db *storm.DB, cleanup func()) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("avisha_test_%d.db", counter.Next()))
	db, err := storm.Open(path)
	if err != nil {
		t.Fatalf("loading database: %v", err)
	}
	return db, func() {
		if err := db.Close(); err != nil {
			t.Logf("closing database: %v", err)
		}
		if err := os.Remove(path); err != nil {
			t.Logf("removing test database: %v", err)
		}
	}
}

// TestPayBill ensures bills can be paid out-of-order, where any excess is stored
// as credit to pay down the next invoice for that service.
//
// @Todo
// - expect error when negative payment
// - handle overpay
//
// How to handle time precision.
// by the time the entities are marshalled between the databased, time has changed
// enough such that Lease.Paid != time.Now()
// One solution is to create a mockable time service that I can swap out for a hardcoded
// time.
func TestPayBill(t *testing.T) {
	mockInvoice := func(amount currency.Currency) Invoice {
		return Invoice{
			Bill: amount,
			Balance: Ledger{
				Debits: []Payment{
					{
						Time:   time.Now(),
						Amount: amount,
					},
				},
			},
		}
	}
	type payment struct {
		InvoiceID int
		Payment
	}
	tests := []struct {
		Label    string
		Input    []Invoice
		Want     []Invoice
		Payments []payment
	}{
		{
			Label: "pay bill exact",
			Input: []Invoice{
				mockInvoice(currency.Dollar * 200),
			},
			Payments: []payment{
				{
					InvoiceID: 1,
					Payment: Payment{
						Time:   time.Now(),
						Amount: currency.Dollar * 200,
					},
				},
			},
			Want: []Invoice{
				{
					ID:   1,
					Bill: currency.Dollar * 200,
					Balance: Ledger{
						Debits: []Payment{
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 200,
							},
						},
						Credits: []Payment{{
							Time:   time.Now(),
							Amount: currency.Dollar * 200,
						}},
					},
					Paid: time.Now(),
				},
			},
		},
		{
			Label: "pay bill in chunks",
			Input: []Invoice{
				mockInvoice(currency.Dollar * 200),
			},
			Payments: []payment{
				{
					InvoiceID: 1,
					Payment: Payment{
						Time:   time.Now(),
						Amount: currency.Dollar * 50,
					},
				},
				{
					InvoiceID: 1,
					Payment: Payment{
						Time:   time.Now(),
						Amount: currency.Dollar * 50,
					},
				},
				{
					InvoiceID: 1,
					Payment: Payment{
						Time:   time.Now(),
						Amount: currency.Dollar * 50,
					},
				},
				{
					InvoiceID: 1,
					Payment: Payment{
						Time:   time.Now(),
						Amount: currency.Dollar * 50,
					},
				},
			},
			Want: []Invoice{
				{
					ID:   1,
					Bill: currency.Dollar * 200,
					Balance: Ledger{
						Debits: []Payment{
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 200,
							},
						},
						Credits: []Payment{
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 50,
							},
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 50,
							},
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 50,
							},
							{
								Time:   time.Now(),
								Amount: currency.Dollar * 50,
							},
						},
					},
					Paid: time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Label, func(t *testing.T) {
			db, cleanup := OpenDB(t)
			defer cleanup()
			for _, invoice := range tt.Input {
				if err := db.Save(&invoice); err != nil {
					t.Fatalf("loading invoice: %v", err)
				}
			}
			app := App{DB: db}
			for _, p := range tt.Payments {
				if err := app.Pay(p.InvoiceID, p.Payment); err != nil {
					t.Fatalf("payment failed: %v", err)
				}
			}
			var (
				got []Invoice
			)
			if err := app.All(&got); err != nil {
				t.Fatalf("loading invoices: %v", err)
			}
			if diff := deep.Equal(got, tt.Want); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
