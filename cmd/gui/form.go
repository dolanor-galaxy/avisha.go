package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

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
