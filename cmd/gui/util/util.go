package util

import (
	"fmt"
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
