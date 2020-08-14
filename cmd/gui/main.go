package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Avisha - Property Management")

	w.CenterOnScreen()
	w.Resize(fyne.Size{Width: 800, Height: 600})

	// TODO:
	// - Setup App and it's services.
	// - Create GUI that consumes app.

	w.SetContent(
		widget.NewVBox(
			NewForm("Tenant Name", "Site Number", "three"),
		),
	)

	w.Show()
	a.Run()
}

// Form with dynamic fields.
type Form struct {
	*widget.Box
	Fields map[string]*Field
	Submit *widget.Button
}

// NewForm allocates a Form with the given field names.
func NewForm(fields ...string) *Form {
	form := &Form{
		Box:    widget.NewVBox(),
		Fields: map[string]*Field{},
		Submit: widget.NewButton("Submit", func() {}),
	}

	for _, field := range fields {
		form.Fields[field] = NewField(field, Vertical)
		form.Box.Append(form.Fields[field])
	}

	form.Box.Append(form.Submit)

	return form
}

// Field widget wraps a label and an entry together.
type Field struct {
	*widget.Box
	Label *widget.Label
	Entry *widget.Entry
}

// Direction of widget stack for the Field.
type Direction = int

const (
	// Vertical places label above entry.
	Vertical = iota
	// Horizontal places label next to entry.
	Horizontal
)

// NewField allocates a Field widget.
func NewField(name string, direction Direction) *Field {
	var (
		container *widget.Box
		label     = widget.NewLabel(name)
		entry     = widget.NewEntry()
	)

	switch direction {
	case Vertical:
		container = widget.NewVBox()
	case Horizontal:
		container = widget.NewHBox()
	default:
		container = widget.NewVBox()
	}

	container.Append(label)
	container.Append(entry)

	return &Field{
		Box:   container,
		Label: label,
		Entry: entry,
	}
}
