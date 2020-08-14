package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	avisha "github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"
)

func init() {
	// Note: Setup spew utility global config.
	// Useful for debugging, subject to change.
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.Indent = "\t"
	spew.Config.SortKeys = true
}

func main() {
	app := avisha.App{
		Storer: &storage.File{
			Path: "db.json",
			Types: map[string]storage.Type{
				"tenant": &TenantType{},
				"site":   &SiteType{},
				"lease":  &LeaseType{},
			},
			Buckets: make(map[string][]json.RawMessage),
		},
		Notifier: &notify.Console{},
	}

	mux := Command{
		Children: []Command{
			tenant,
		},
	}

	stdin := bufio.NewReader(os.Stdin)

	for {
		line, _ := stdin.ReadString('\n')
		if err := mux.Handle(&app, strings.Fields(line)); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			spew.Dump(app.Storer)
		}
	}
}

type TenantType struct{}

func (t *TenantType) New() interface{} {
	return &avisha.Tenant{}
}

func (t *TenantType) Is(v interface{}) bool {
	_, ok := v.(avisha.Tenant)
	return ok
}

type SiteType struct{}

func (t *SiteType) New() interface{} {
	return &avisha.Site{}
}

func (t *SiteType) Is(v interface{}) bool {
	_, ok := v.(avisha.Site)
	return ok
}

type LeaseType struct{}

func (t *LeaseType) New() interface{} {
	return &avisha.Lease{}
}

func (t *LeaseType) Is(v interface{}) bool {
	_, ok := v.(avisha.Lease)
	return ok
}
