package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/currency"
)

// LoadFakeData loads in some pre-baked data for testing and development purposes.
func LoadFakeData(db *storm.DB) error {
	randomAddress := func() avisha.Address {
		var streets = []string{
			"Riverhead Road",
			"Sacovia Place",
			"Lightning Crescent",
			"Money Hill",
		}
		return avisha.Address{
			Unit:   rand.Intn(20),
			Number: rand.Intn(256),
			Street: streets[rand.Intn(len(streets)-1)],
			City:   "Elysium",
		}
	}
	var (
		tenantID = 0
	)
	makeTenant := func(name string) avisha.Tenant {
		var (
			handle string
			domain = "example"
		)
		for ii, field := range strings.Fields(name) {
			if ii == 0 {
				handle = field
			} else if ii == 1 {
				domain = field
			}
		}
		tenantID++
		return avisha.Tenant{
			ID:      tenantID,
			Name:    name,
			Contact: fmt.Sprintf("%s@%s.com", handle, domain),
			Address: randomAddress(),
		}
	}
	sites := []avisha.Site{
		{
			ID:       1,
			Number:   "1",
			Dwelling: avisha.Cabin,
		},
		{
			ID:       2,
			Number:   "2",
			Dwelling: avisha.Flat,
		},
		{
			ID:       3,
			Number:   "3",
			Dwelling: avisha.House,
		},
		{
			ID:       4,
			Number:   "4",
			Dwelling: avisha.Cabin,
		},
	}
	tenants := []avisha.Tenant{
		makeTenant("Sheriff Hardin"),
		makeTenant("Phantom Menace"),
		makeTenant("Jack Mordaunt"),
		makeTenant("Tony Stark"),
	}
	settings := avisha.Settings{
		Landlord: avisha.Landlord{
			Name:  "FooInc",
			Email: "admin@fooinc.com",
			Phone: "123 456 789",
			Address: avisha.Address{
				Unit:   12,
				Number: 128,
				Street: "Foo Crescent",
				City:   "Elysium",
			},
		},
		Bank: avisha.Bank{
			Name:    "Bank of Elysium",
			Account: "123345567",
		},
		Defaults: avisha.Defaults{
			UnitCost:   currency.Dollar,
			RentCycle:  14 * 24 * time.Hour,
			InvoiceNet: 14 * 24 * time.Hour,
			GST:        10,
		},
	}
	leases := []avisha.Lease{}
	for ii, site := range sites {
		ii := ii
		tenant := tenants[ii]
		leases = append(leases, avisha.Lease{
			Site:   site.ID,
			Tenant: tenant.ID,
			Term: avisha.Term{
				Start:    time.Now(),
				Duration: time.Duration(rand.Intn(365+50)) * 24 * time.Hour,
			},
			Rent: currency.Dollar * currency.Currency(rand.Intn(200+50)),
		})
	}
	for _, s := range sites {
		if err := db.Save(&s); err != nil {
			return err
		}
	}
	for _, t := range tenants {
		if err := db.Save(&t); err != nil {
			return err
		}
	}
	for _, l := range leases {
		if err := db.Save(&l); err != nil {
			return err
		}
	}
	if err := db.Set("settings", "global", &settings); err != nil {
		return err
	}
	return nil
}
