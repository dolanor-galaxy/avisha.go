package util

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"strconv"
	"strings"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/currency"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Rect struct {
	Color color.NRGBA
	Size  image.Point
	Radii unit.Value
}

func (r Rect) Layout(gtx C) D {
	return DrawRect(gtx, r.Color, r.Size, r.Radii)
}

// DrawRect creates a rectangle of the provided background color with
// Dimensions specified by size and a corner radius (on all corners)
// specified by radii.
func DrawRect(gtx C, background color.NRGBA, size image.Point, radii unit.Value) D {
	defer op.Push(gtx.Ops).Pop()
	rr := float32(gtx.Px(radii))
	clip.Rect{Max: size}.Add(gtx.Ops)
	paint.ColorOp{
		Color: background,
	}.Add(gtx.Ops)
	if rr != 0 {
		clip.RRect{
			Rect: f32.Rectangle{
				Max: layout.FPt(size),
			},
			NW: rr,
			NE: rr,
			SE: rr,
			SW: rr,
		}.Add(gtx.Ops)
	}
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// ParseDate parses a time object from a textual dd/mm/yyyy format.
func ParseDate(s string) (date time.Time, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 3 {
		return date, fmt.Errorf("must be dd/mm/yyyy")
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return date, fmt.Errorf("year not a number: %s", parts[2])
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return date, fmt.Errorf("month not a number: %s", parts[2])
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return date, fmt.Errorf("day not a number: %s", parts[2])
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// @Todo Consider package api for these form utility functions.

// ParseInt parses an integer from digit characters.
func ParseInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	return n, nil
}

// ParseCurrency parses a dollars from digit characters.
func ParseCurrency(s string) (currency.Currency, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	return currency.Currency(n) * currency.Dollar, nil
}

// ParseInt parses an unsigned integer from digit characters.
func ParseUint(s string) (uint, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	} else if n < 1 {
		return 0, fmt.Errorf("must be an amount greater than 0")
	}
	return uint(n), nil
}

// ParseDay parses a day from digit characters.
func ParseDay(s string) (time.Duration, error) {
	n, err := ParseUint(s)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Hour * 24 * time.Duration(n), nil
}

// FieldRequired ensures that a string is not empty.
func FieldRequired(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", fmt.Errorf("required")
	}
	return s, nil
}

// FormatTime formats a time object into a string.
func FormatTime(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Day(), t.Month(), t.Year())
}

// FlexStrategy renders flexed with the given weight if the axis matches.
// Otherwise the widget is rendered rigid and weight is ignored.
func FlexStrategy(weight float32, flex, actual layout.Axis, w layout.Widget) layout.FlexChild {
	if flex == actual {
		return layout.Flexed(weight, w)
	}
	return layout.Rigid(w)
}

// UtilityInvoiceDocument renders utility invoices to an html document.
type UtilityInvoiceDocument struct {
	Invoice avisha.UtilityInvoice
	Lease   avisha.Lease
	Tenant  avisha.Tenant
	Site    avisha.Site
}

// Render the document into a buffer.
func (doc UtilityInvoiceDocument) Render() (*bytes.Buffer, error) {
	tmpl, err := template.
		New("utility-invoice-document").
		Parse(strings.TrimSpace(UtilityInvoiceTemplateLiteral))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}
	by := new(bytes.Buffer)
	if err := tmpl.Execute(by, doc); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return by, nil
}

var UtilityInvoiceTemplateLiteral = `
<html>
	<head>
		<title>Invoice {{.Invoice.ID}}</title>
	</head>
	<body>
		<h1>Invoice {{.Invoice.ID}}</h1>
	</body>
</html>
`
