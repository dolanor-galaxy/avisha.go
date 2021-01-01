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

// ParseFloat parses a floating point number from digit characters.
func ParseFloat(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	return n, nil
}

// ParseCurrency parses a dollars from digit characters.
func ParseCurrency(s string) (c currency.Currency, err error) {
	s = strings.TrimPrefix(s, "$")
	parts := strings.Split(s, ".")
	d, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("must be a valid number")
	}
	c += (currency.Currency(d) * currency.Dollar)
	if len(parts) > 1 {
		var (
			fraction = parts[1]
			length   = len(fraction)
		)
		// Ignore extra digits.
		if length > 4 {
			length = 4
		}
		// Pad with zeros for correct precision.
		if length < 4 {
			fraction += strings.Repeat("0", 4-length)
		}
		mills, err := strconv.Atoi(fraction)
		if err != nil {
			return 0, fmt.Errorf("must be a valid number")
		}
		c += (currency.Currency(mills) * currency.Mill)
	}
	return c, nil
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
	// @Todo drive previous reading from history
	// @Todo use history to display unpaid invoices
	History  []*avisha.UtilityInvoice
	Previous avisha.UtilityInvoice
	Invoice  avisha.UtilityInvoice

	Lease    avisha.Lease
	Tenant   avisha.Tenant
	Site     avisha.Site
	Settings avisha.Settings
}

// Render the document into a buffer.
func (doc UtilityInvoiceDocument) Render() (*bytes.Buffer, error) {
	tmpl, err := template.
		New("utility-invoice-document").
		Funcs(template.FuncMap{
			"date": func(t time.Time) string {
				return t.Format("Monday, 2 January 2006")
			},
			"generateReference": func() string {
				// Get the first three letters of the last name.
				fields := strings.Fields(doc.Tenant.Name)
				name := strings.ToUpper(fields[len(fields)-1][0:3])
				site := strings.ToUpper(doc.Site.Number)
				return fmt.Sprintf("%s-S.%s POWR", name, site)
			},
		}).
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

// UtilityInvoiceTemplateLiteral contains the literal html used to generate
// an html invoice, which can be saved as pdf by most browers.
//
// @Todo render direct to pdf.
var UtilityInvoiceTemplateLiteral = `
<!doctype html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0"> 
		<link rel="stylesheet" href="https://vanillacss.com/vanilla.css" media="all">
		<title>Invoice {{.Invoice.ID}}</title>
		<style>
			body{
				margin: 0 auto;
				max-width: 50rem;
			}
			@media(max-width: 50rem) {
				body {
					padding: 10px;
				}
			}
			table,tbody {
				text-align: center;
			}
			table td {
				padding: 0.25rem;
			}
			blockquote p:last-child {
				margin-bottom: 0;
			}
			cards {
				display: flex;
				flex-direction: row;
				justify-content: space-between;
			}
			card {
				width: 100%;
				margin: 0.5rem;
				border: 1px solid var(--primary-color) !important;
			}
			card.no-border {
				border: 0px;
			}
			card header {
				width: 100%;
				padding: 0 1rem;
				font-weight: bold;
				background: var(--secondary-color);
				border-bottom: 1px solid var(--primary-color) !important;
			}
			card p  {
				width: 100%;
				padding: 1rem;
				margin: 0;
			}
			.compact {
				margin: 0;
				padding: 0;
			}
			.compact li {
				margin: 0;
				padding: 0;
				margin-left: 4rem;
			}
			table caption {
				margin: 0;
				padding: 0.25rem;
				text-align: left;
				font-weight: bold;
			}
			table tbody {
				text-align: center !important;
			}
			td var {
				font-weight: normal !important;
			}
			@media print {
				body {
					font-size: 14pt;
				}
				// Printed page already has top margin.
				article:first-of-type h1 {
					margin-top: 0;
				}
				table {
					page-break-inside: avoid;
					margin: 1rem 0;
				}
				card {
					page-break-inside: avoid;
				}
			}
		</style>
	</head>
	<body id="top" role="document">
		<article id="preamble">
			<header><h1>Tax Invoice / Statement {{.Invoice.ID}}</h1></header>
			<cards>
				<card class="no-border">
					<p>
						<b>AVISHA GROUP LTD</b>
						</br>
						Property Management Services
						</br>
						GST No. 125544207
						</br>
						{{.Settings.Landlord.Address}}
					</p>
				</card>
				<card class="no-border">
					<p>
						Statement Date
						</br>
						{{date .Invoice.Issued}}
					</p>
				</card>
			</cards>
			<cards>
				<card>
					<header>Site</header>
					<p>
						Number: {{.Site.Number}}
						</br>
						Type: {{.Site.Dwelling}}
						</br>
						Period: <var>{{.Invoice.Period}}</var>
						</br>
						Service: <b>Electricity</b>
					</p>
				</card>
				<!-- @Todo polymorph the service? -->
				<card>
					<header>Bill To</header>
					<p>
						{{.Tenant.Name}}
						</br>
						{{.Tenant.Address}}
						</br>
						{{.Tenant.Contact}}
					</p>
				</card>
			</cards>
		</article>
		<article id="activity">
			<header><h1>Activity</h1></header>
			<table>
				<caption>Previous Activity</caption>
				<thead>
					<tr>
						<th>Invoice</th>
						<th>Bill</th>
						<th>Received</th>
						<th>Outstanding</th>
					</tr>
				</thead>
				<tbody>
					{{range $invoice := .History}}
						{{if not $invoice.IsPaid}}
						<tr>
							<td><var>{{$invoice.ID}}</var></td>
							<td><var>{{$invoice.Bill}}</var></td>
							<!-- When was the most recent payment received, if at all? --> 
							<td><var></var></td>
							<td><var>{{$invoice.Balance}}</var></td>
						</tr>
						{{end}}
					{{end}}
				</tbody>
			</table>
			<table>
				<caption>Current Activity</caption>
				<thead>
					<tr>
						<th>Unit Cost</th>
						<th>Previous Reading</th>
						<th>Current Reading</th>
						<th>Units Used</th>
						<th>Activity Charge</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td><var>{{.Invoice.UnitCost}}</var></td>
						<td><var>{{.Previous.Reading}}</var></td>
						<td><var>{{.Invoice.Reading}}</var></td>
						<td><var>{{.Invoice.UnitsConsumed}}</var></td>
						<!-- @Todo utilities cost, not total bill -->
						<td><var>{{.Invoice.Charges.Activity}}</var></td>
					</tr>
				</tbody>
			</table>
			<table>
				<caption>Charges</caption>
				<thead>
					<tr>
						<th>Line Charge</th>
						<th>Late Fee</th>
						<!-- @Todo pull gst from settings -->
						<th>GST ({{.Invoice.GST}}%)</th>
						<th>Total Charges</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td><var>{{.Invoice.Charges.LineCharge}}</var></td>
						<td><var>{{.Invoice.Charges.LateFee}}</var></td>
						<td><var>{{.Invoice.Charges.GST}}</var></td>
						<td><var>{{.Invoice.Bill}}</var></td>
					</tr>
				</tbody>
			</table>
			<blockquote>
				<p>
					Total Amount Due by <time>{{date .Invoice.Due}}</time> <var>{{.Invoice.Bill}}</var>
					</br>
					<small>(please note late payment fee will be charged if payment not received by due date)</small>
				</p>
			</blockquote>
			<cards>
				<card>
					<header>Make Payable To</header>
					<p>
						<b>Bank Acc:</b> {{.Settings.Bank.Name}} <var>{{.Settings.Bank.Account}}</var>
						</br>
						<b>Reference:</b> {{generateReference}}
						</br>
						<!-- @Todo list contact items in generic fashion -->
						<b>Email:</b> {{.Settings.Landlord.Email}}
						</br>
						<b>Phone:</b> {{.Settings.Landlord.Phone}}
					</p>
				</card>
			</cards>
		</article>
	</body>
</html>
`
