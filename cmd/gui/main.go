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
	w.Resize(fyne.Size{Width: 800, Height: 600})
	// w.SetPadded(true)

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

	// main := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), sidebar, content)

	main := widget.NewHBox(sidebar, content)

	// TODO:
	// - Setup App and it's services.
	// - Create GUI that consumes app.
	w.SetContent(main)

	w.Show()
	a.Run()
}

// Form renders a form with the container fields.
type Form struct {
	Title  string
	Fields []Field

	OnSubmit func(form *Form)
	OnCancel func(form *Form)
	// Validate func(form *Form) error
}

// Data returns fieldwise data "field:value".
// All field data is encoded as a string for the purposes of the form.
func (form *Form) Data() map[string]string {
	m := map[string]string{}
	for _, field := range form.Fields {
		m[field.Name] = field.Entry.Text
	}
	return m
}

// Build a form widget.
func (form *Form) Build() fyne.Widget {
	box := widget.NewVBox()
	buttons := widget.NewHBox()

	for ii, field := range form.Fields {
		ii := ii

		fieldBox := widget.NewVBox()

		switch field.Kind {
		case Text:
			form.Fields[ii].Entry = widget.NewEntry()
		default:
			form.Fields[ii].Entry = widget.NewEntry()
		}

		fieldBox.Append(field.Label)
		fieldBox.Append(form.Fields[ii].Entry)

		box.Append(fieldBox)
	}

	buttons.Append(widget.NewButton("Submit", func() {
		if form.OnSubmit != nil {
			form.OnSubmit(form)
		}
		for _, field := range form.Fields {
			field.Entry.SetText("")
		}
	}))

	buttons.Append(widget.NewButton("Cancel", func() {
		if form.OnCancel != nil {
			form.OnCancel(form)
		}
		for _, field := range form.Fields {
			field.Entry.SetText("")
		}
	}))

	box.Append(buttons)

	return box
}

// Field associates a label with an input and validation rules.
type Field struct {
	Name  string
	Entry *widget.Entry
	Label *widget.Label
	Kind  InputKind
	// Validate func(field *Field) error
}

// InputKind specifies the type of data input desired.
type InputKind int

const (
	// Text allows plain text.
	Text InputKind = iota
	// Date allows date notation.
	Date
	// Number allows numeric characters.
	Number
)
