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
	)
	if l.Balance < 0 {
		sign = "-"
	}
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return material.Label(l.Th.Primary(), l.Th.TextSize, l.Name).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			lb := material.Label(l.Th.Danger(), l.Th.TextSize, fmt.Sprintf(" %s$%.2f", sign, balance))
			lb.Color = l.Th.Success().Color.Primary
			if balance < 0 {
				lb.Color = l.Th.Danger().Color.Primary
			}
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
