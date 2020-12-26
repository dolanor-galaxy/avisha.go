package views

import (
	"gioui.org/layout"
	"git.sr.ht/~whereswaldon/materials"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
)

// AddressForm for manipulating structured address values.
type AddressForm struct {
	Address *avisha.Address
	Fields  struct {
		Unit   materials.TextField
		Number materials.TextField
		Street materials.TextField
		City   materials.TextField
	}
	Form widget.Form
}

func (f *AddressForm) Submit() (address avisha.Address, ok bool) {
	return *f.Address, f.Form.Submit()
}

func (f *AddressForm) Clear() {
	f.Form.Clear()
}

func (f *AddressForm) Load(address *avisha.Address) {
	f.Address = address
	f.Form.Load([]widget.Field{
		{
			Value: widget.IntValuer{Value: &f.Address.Unit},
			Input: &f.Fields.Unit,
		},
		{
			Value: widget.RequiredValuer{Valuer: widget.IntValuer{Value: &f.Address.Number}},
			Input: &f.Fields.Number,
		},
		{
			Value: widget.RequiredValuer{Valuer: widget.TextValuer{Value: &f.Address.Street}},
			Input: &f.Fields.Street,
		},
		{
			Value: widget.RequiredValuer{Valuer: widget.TextValuer{Value: &f.Address.City}},
			Input: &f.Fields.City,
		},
	})
}

func (f *AddressForm) Layout(gtx C, th *style.Theme) D {
	f.Form.Validate(gtx)
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return f.Fields.Unit.Layout(gtx, th.Dark(), "Unit")
		}),
		layout.Rigid(func(gtx C) D {
			return f.Fields.Number.Layout(gtx, th.Dark(), "Number")
		}),
		layout.Rigid(func(gtx C) D {
			return f.Fields.Street.Layout(gtx, th.Dark(), "Street")
		}),
		layout.Rigid(func(gtx C) D {
			return f.Fields.City.Layout(gtx, th.Dark(), "City")
		}),
	)
}
