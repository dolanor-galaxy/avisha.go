package main

import (
	"sync"

	"gioui.org/layout"
)

// Route between named views.
type Router struct {
	sync.Mutex
	Static func(gtx Ctx, r *Router) Dims
	Routes map[string]Route
	Stack  []string
}

func (r *Router) Pop() {
	defer r.lock()()
	if len(r.Stack) > 1 {
		r.Stack = r.Stack[:len(r.Stack)-1]
	}
}

func (r *Router) Push(s string) {
	defer r.lock()()
	if r.Stack[len(r.Stack)-1] != s {
		if _, ok := r.Routes[s]; ok {
			r.Stack = append(r.Stack, s)
		}
	}
}

func (r *Router) Update(gtx Ctx) {
	if name, ok := r.Active().Route(); ok {
		if name == RouteBack {
			r.Pop()
		} else if name != "" {
			r.Push(name)
		}
	}
}

// Layout static content as rigid, then layout the active route.
func (r *Router) Layout(gtx Ctx) Dims {
	r.Update(gtx)
	// return r.Active().Layout(gtx)
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx Ctx) Dims {
			return r.Static(gtx, r)
		}),
		layout.Flexed(1, func(gtx Ctx) Dims {
			return r.Active().Layout(gtx)
		}),
	)
}

func (r *Router) Name() string {
	return r.Stack[len(r.Stack)-1]
}

func (r *Router) Active() Route {
	return r.Routes[r.Stack[len(r.Stack)-1]]
}

func (r *Router) lock() func() {
	r.Lock()
	return func() { r.Unlock() }
}

type Route interface {
	Route() (string, bool)
	Actions() []layout.Widget
	Layout(gtx Ctx) Dims
}

const RouteBack = "back"
