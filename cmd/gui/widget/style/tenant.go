package style

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha.go"
)

// TenantLabelStyle renders a Tenant label.
type TenantLabelStyle struct {
	Tenant avisha.Tenant
	Theme  *material.Theme

	NameColor color.NRGBA
	NameSize  unit.Value

	IDColor color.NRGBA
	IDSize  unit.Value
}

func (t TenantLabelStyle) Layout(gtx C) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			l := material.Body1(t.Theme, t.Tenant.Name)
			l.Color = t.NameColor
			l.TextSize = t.NameSize
			return l.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return D{Size: image.Point{X: gtx.Px(unit.Dp(2))}}
		}),
		layout.Rigid(func(gtx C) D {
			l := material.Body1(t.Theme, fmt.Sprintf("#%d", t.Tenant.ID))
			l.Color = t.IDColor
			l.TextSize = t.IDSize
			return l.Layout(gtx)
		}),
	)
}

// TenantLabel renders a Tenant in a consistent way.
// Tenant name combined with a de-emphasized identifier.
func TenantLabel(th *material.Theme, t avisha.Tenant) TenantLabelStyle {
	return TenantLabelStyle{
		Tenant:    t,
		Theme:     th,
		NameColor: th.Fg,
		NameSize:  unit.Dp(20),
		IDColor:   WithAlpha(th.Fg, 100),
		IDSize:    unit.Dp(15),
	}
}
