package main

import (
	"bufio"
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

// tenant register name=jack contact=jack@protonmail.com

func main() {
	app := avisha.App{
		Storer: storage.FileStorage("target/cli/db.json").
			With(&avisha.Tenant{}).
			With(&avisha.Site{}).
			With(&avisha.Lease{}).
			MustLoad(),
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
