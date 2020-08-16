package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"
)

func main() {
	appo := avisha.App{
		Storer: storage.FileStorage("target/cli/db.json").
			With(&avisha.Tenant{}).
			With(&avisha.Site{}).
			With(&avisha.Lease{}).
			MustLoad(),
		Notifier: &notify.Console{},
	}

	a := app.NewWithID("avisha")
	w := a.NewWindow("Avisha - Property Management")

	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 1200, Height: 800})

	sidebar := widget.NewVBox()
	content := widget.NewVScrollContainer(
		fyne.NewContainerWithLayout(layout.NewCenterLayout(), (&Form{
			Title: "Lease",
			Fields: []Field{
				{
					Name:  "tenant",
					Label: widget.NewLabel("Tenant Name"),
				},
				{
					Name:  "site",
					Label: widget.NewLabel("Site Number"),
				},
				{
					Name:  "start",
					Label: widget.NewLabel("Start"),
				},
				{
					Name:  "duration",
					Label: widget.NewLabel("Duration (days)"),
				},
			},
			OnSubmit: func(form *Form) {
				data := form.Data()

				var (
					tenant string
					site   string
					term   avisha.Term
					rent   uint
				)

				if t, ok := data["tenant"]; ok {
					tenant = t
				}

				if s, ok := data["site"]; ok {
					site = s
				}

				if start, ok := data["start"]; ok {
					if date, err := time.Parse(time.RFC822, start); err != nil {
						term.Start = date
					}
				}

				if duration, ok := data["duration"]; ok {
					if duration, err := strconv.Atoi(duration); err != nil {
						term.Days = duration
					}
				}

				if err := appo.CreateLease(tenant, site, term, rent); err != nil {
					// TODO: Display error to user via GUI
					fmt.Printf("error: creating lease: %s\n", err)
				}
			},
		}).Build()),
	)

	main := widget.NewHBox(sidebar, content)

	w.SetContent(main)

	w.Show()
	a.Run()
}
