package main

import (
	"time"

	"gioui.org/layout"
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

func (p *Page) Route() (string, bool) {
	defer func() { p.reroute = "" }()
	if ok := p.reroute != ""; ok {
		return p.reroute, ok
	}
	return "", false
}
