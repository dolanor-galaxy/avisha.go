package widget

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackmordaunt/avisha.go/cmd/gui/util"
	"github.com/jackmordaunt/avisha.go/currency"
)

// Valuer implements a bi-directional mapping between textual data and structured
// data.
// Valuer contains the precise validation logic, which is expresssed by the
// error return.
// @Taxonomy Is there a better name for this? FieldMapper, Mapper, Value.
type Valuer interface {
	// To converts the value into a string.
	To() (text string, err error)
	// From parses the value from a string.
	From(text string) (err error)
	// Clear resets the value.
	Clear()
}

// Input can handle text content and error content.
// Currently based on the `materials.TextField` method set.
type Input interface {
	// @Cleanup Could Text and SetText be collapsed into a single method?
	Text() string
	SetText(string)
	// @Cleanup Could SetError and ClearError be collapsed into a single method?
	SetError(string)
	ClearError()
	// @Note We want some way of validating on input event, but we don't care
	// about the representation of the event. For now we will hard code to the
	// widget.EditorEvent, but maybe that's not appropriate.
	Events() []EditorEvent
}

// Field associates a value with a name.
// Name is the formatted title of the field, suitable for rendering to the UI.
type Field struct {
	Value Valuer
	Input Input
}

// Validate the field by running the text through the Valuer.
// Precise validation logic is implemented by the Valuer.
// Returns a boolean indication success.
func (field *Field) Validate() bool {
	err := field.Value.From(field.Input.Text())
	if err != nil {
		field.Input.SetError(err.Error())
	} else {
		field.Input.ClearError()
	}
	return err == nil
}

// Form manipulates fields in a consistent way.
type Form struct {
	Fields []Field
}

// Load values into inputs.
func (f *Form) Load(fields []Field) {
	if len(fields) > 0 {
		f.Fields = fields
	}
	for _, field := range f.Fields {
		if text, err := field.Value.To(); err != nil {
			field.Input.SetError(err.Error())
		} else {
			field.Input.ClearError()
			field.Input.SetText(text)
		}
	}
}

// Submit batch validates the fields and returns a boolean indication success.
// If true, all the fields validated and you can use the data.
func (f *Form) Submit() (ok bool) {
	ok = true
	for _, field := range f.Fields {
		if !field.Validate() {
			ok = false
		}
	}
	return ok
}

// Validate form fields in realtime.
func (f *Form) Validate(gtx C) {
	for _, field := range f.Fields {
		for range field.Input.Events() {
			field.Validate()
		}
	}
}

func (f *Form) Clear() {
	for _, field := range f.Fields {
		field.Value.Clear()
		text, _ := field.Value.To()
		field.Input.ClearError()
		field.Input.SetText(text)
	}
}

// Basic Value implementations.
// @Cleanup Move to appropriate package.

// IntValuer maps integers to text.
type IntValuer struct {
	Value *int
}

func (v IntValuer) To() (string, error) {
	return strconv.Itoa(*v.Value), nil
}

func (v IntValuer) From(text string) (err error) {
	*v.Value, err = util.ParseInt(text)
	return err
}

func (v IntValuer) Clear() {
	*v.Value = 0
}

// FloatValuer maps floating points to text.
type FloatValuer struct {
	Value *float64
}

func (v FloatValuer) To() (string, error) {
	return strconv.FormatFloat(*v.Value, 'f', 2, 64), nil
}

func (v FloatValuer) From(text string) (err error) {
	*v.Value, err = util.ParseFloat(text)
	return err
}

func (v FloatValuer) Clear() {
	*v.Value = 0
}

// CurrencyValuer maps currency to text.
type CurrencyValuer struct {
	Value *currency.Currency
}

func (v CurrencyValuer) To() (string, error) {
	return strings.TrimPrefix(v.Value.String(), "$"), nil
}

func (v CurrencyValuer) From(text string) (err error) {
	*v.Value, err = util.ParseCurrency(text)
	return err
}

func (v CurrencyValuer) Clear() {
	*v.Value = 0
}

// TextValuer wraps a text value.
type TextValuer struct {
	Value *string
}

func (v TextValuer) To() (string, error) {
	return *v.Value, nil
}

func (v TextValuer) From(text string) error {
	*v.Value = text
	return nil
}

func (v TextValuer) Clear() {
	*v.Value = ""
}

type DaysValuer struct {
	Value *time.Duration
}

func (v DaysValuer) To() (string, error) {
	days := (*v.Value) / (time.Hour * 24)
	return strconv.Itoa(int(days)), nil
}

func (v DaysValuer) From(text string) (err error) {
	*v.Value, err = util.ParseDay(text)
	return err
}

func (v DaysValuer) Clear() {
	*v.Value = time.Hour * 24
}

type DateValuer struct {
	Value *time.Time
}

func (v DateValuer) To() (string, error) {
	return fmt.Sprintf("%d/%d/%d", v.Value.Day(), v.Value.Month(), v.Value.Year()), nil
}

func (v DateValuer) From(text string) (err error) {
	*v.Value, err = util.ParseDate(text)
	return err
}

func (v DateValuer) Clear() {
	*v.Value = time.Now()
}

// RequiredValuer ensures the field is not empty.
type RequiredValuer struct {
	Valuer
}

func (v RequiredValuer) From(text string) error {
	if len(strings.TrimSpace(text)) == 0 {
		return fmt.Errorf("required")
	}
	return v.Valuer.From(text)
}
