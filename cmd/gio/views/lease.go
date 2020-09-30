package views

import (
	"fmt"
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/storage"
)

type Lease struct {
	*avisha.App
	*material.Theme
	list layout.List
	once sync.Once
}

func (l *Lease) Route() (string, bool) {
	return "", false
}

func (l *Lease) Actions() []layout.Widget {
	return nil
}

func (l *Lease) Layout(gtx Ctx) Dims {
	l.once.Do(func() {
		l.list.Axis = layout.Vertical
		l.list.ScrollToEnd = false
	})
	var (
		leases []*avisha.Lease
	)
	l.App.List(func(ent storage.Entity) bool {
		l, ok := ent.(*avisha.Lease)
		if ok {
			leases = append(leases, l)
		}
		return ok
	})
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx Ctx) Dims {
		return l.list.Layout(gtx, len(leases), func(gtx Ctx, index int) Dims {
			var (
				lease = leases[index]
			)
			return material.Label(
				l.Theme,
				unit.Dp(14),
				fmt.Sprintf("%s - %s: %+v", lease.Site, lease.Tenant, lease.Term)).Layout(gtx)
		})
	})
}
