package main

import (
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Page displays a page with common configuraiton.
type Page struct {

	// Nav state.
	Home  widget.Clickable
	Other widget.Clickable
	Back  widget.Clickable

	reroute string
	pending string
	time    time.Time
}

func (p *Page) Layout(gtx Ctx) Dims {
	p.Update(gtx)
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Flexed(1, func(gtx Ctx) Dims {
			return Dims{}
		}),
	)
}

func (p *Page) Update(gtx Ctx) {
	// Hack(UX): wait for ink animation before re-routing.
	// Alternative would be to have the topbar own the action state,
	// which would allow ink animations intra route.
	route := func(n string) {
		p.pending = n
		p.time = time.Now().Add(250 * time.Millisecond)
	}
	if p.time.After(time.Now()) {
		return
	}
	if p.Home.Clicked() {
		route("home")
	}
	if p.Other.Clicked() {
		route("other")
	}
	if p.Back.Clicked() {
		route(RouteBack)
	}
	if p.time.Before(time.Now()) {
		p.reroute = p.pending
		p.pending = ""
	}
}

func (p *Page) Actions() []layout.Widget {
	return []layout.Widget{
		func(gtx Ctx) Dims {
			return material.Button(th, &p.Home, "Home").Layout(gtx)
		},
		func(gtx Ctx) Dims {
			return material.Button(th, &p.Other, "Other").Layout(gtx)
		},
		func(gtx Ctx) Dims {
			return material.Button(th, &p.Back, "Go Back").Layout(gtx)
		},
	}
}

func (p *Page) ReRoute() (string, bool) {
	defer func() { p.reroute = "" }()
	if ok := p.reroute != ""; ok {
		return p.reroute, ok
	}
	return "", false
}

// TopBar renders a title and the provided actions.
type TopBar struct {
}

func (bar TopBar) Layout(gtx Ctx, r *Router) Dims {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx Ctx) Dims {
			return Rect{
				Color: color.RGBA{B: 100, R: 50, A: 255},
				Size:  gtx.Constraints.Max,
			}.Layout(gtx)
		}),
		layout.Stacked(func(gtx Ctx) Dims {
			actions := r.Active().Actions()
			items := []layout.FlexChild{
				layout.Flexed(float32(len(actions)+1), func(gtx Ctx) Dims {
					th.Color.Text = color.RGBA{R: 255, G: 255, B: 255, A: 255}
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx Ctx) Dims {
						return material.Label(th, unit.Dp(24), r.Name()).Layout(gtx)
					})
				}),
			}
			for _, action := range actions {
				action := action
				items = append(items, layout.Flexed(1, func(gtx Ctx) Dims {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, action)
				}))
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, items...)
		}),
	)
}
