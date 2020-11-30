package style

import (
	"fmt"
	"image"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"git.sr.ht/~whereswaldon/materials"
)

// DateInput allows inputting of date information textually.
// Composed of three inputs: day, month, year.
type DateInput struct {
	Day   materials.TextField
	Month materials.TextField
	Year  materials.TextField
}

func (input *DateInput) Set(date time.Time) {
	var (
		day, month, year = "", "", ""
	)
	if date != (time.Time{}) {
		day = fmt.Sprintf("%d", date.Day())
		month = fmt.Sprintf("%d", date.Month())
		year = fmt.Sprintf("%d", date.Year())
	}
	input.Day.SetText(day)
	input.Month.SetText(month)
	input.Year.SetText(year)
}

func (input *DateInput) Date() (time.Time, error) {
	year, err := strconv.Atoi(input.Year.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("year: not a number")
	}
	if year < 0 {
		return time.Time{}, fmt.Errorf("year: out of bounds (must be positive number) got %d", year)
	}
	month, err := strconv.Atoi(input.Month.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("month: not a number")
	}
	if month < int(time.January) || month > int(time.December) {
		return time.Time{}, fmt.Errorf("month: out of bounds (1-12) got %d", month)
	}
	day, err := strconv.Atoi(input.Day.Text())
	if err != nil {
		return time.Time{}, fmt.Errorf("day: not a number")
	}
	if day < 1 || day > 31 {
		return time.Time{}, fmt.Errorf("day: out of bounds (1-31) got %d", day)
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

func (input *DateInput) Layout(gtx C, th *Theme) D {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(
		gtx,
		layout.Flexed(1, func(gtx C) D {
			return input.Day.Layout(gtx, th.Primary(), "Day")
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
		}),
		layout.Flexed(1, func(gtx C) D {
			return input.Month.Layout(gtx, th.Primary(), "Month")
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(10))}}
		}),
		layout.Flexed(1, func(gtx C) D {
			return input.Year.Layout(gtx, th.Primary(), "Year")
		}),
	)
}
