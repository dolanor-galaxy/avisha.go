package style

import (
	"fmt"
	"math"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

// ServiceLabelStyle renders the balance of a named service.
type ServiceLabelStyle struct {
	Th      *Theme
	Name    string
	Balance float64
}

func (l *ServiceLabelStyle) Layout(gtx C) D {
	var (
		balance = math.Abs(l.Balance)
		sign    = ""
		color   = l.Th.Success().Fg
	)
	if l.Balance < 0 {
		sign = "-"
		color = l.Th.Danger().Fg
	}
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return material.Label(l.Th.Dark(), l.Th.TextSize, l.Name).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			lb := material.Label(l.Th.Dark(), l.Th.TextSize, fmt.Sprintf(" %s$%.2f", sign, balance))
			lb.Color = color
			return lb.Layout(gtx)
		}),
	)
}

// ServiceLabel renders the balance of a named service.
func ServiceLabel(th *Theme, name string, balance float64) *ServiceLabelStyle {
	return &ServiceLabelStyle{
		Th:      th,
		Name:    name,
		Balance: balance,
	}
}
