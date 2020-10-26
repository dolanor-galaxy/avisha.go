package nav

import (
	"sync"

	"gioui.org/layout"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// Router implements routing between named views.
// Views can implement optional interfaces to allow re-routing, data passing
// and provide context to views that exist outside the Router.
// Static content is rendered independent of the current Route.
type Router struct {
	sync.Mutex
	Routes map[string]View
	Stack  []string
}

// View types render themselves.
type View interface {
	Layout(gtx C) D
}

// Titled views have pretty names.
type Titled interface {
	Title() string
}

// ReRouter views can signal to reroute to another Route by name.
// Arbitrary data can be passed to the target route.
// The ReRouter is responisble for passing the right data to the right route.
//
// Note: This creates an implicit dependency on the both route name, and the
// route data. There may be better ways to handle coupling between routes.
type ReRouter interface {
	ReRoute() (name string, data interface{})
}

// Contexter views provide pieces of UI to be rendered externally.
type Contexter interface {
	Context() []layout.Widget
}

// Receiver views accept arbitrary data on re-route.
// Typically a receiver will type switch on the data it cares about.
type Receiver interface {
	Receive(data interface{})
}

// RouteBack is a special route that tells the Router to route
// to the previous view.
const RouteBack = "back"

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

func (r *Router) Update(gtx C) {
	if rerouter, ok := r.Active().(ReRouter); ok {
		to, data := rerouter.ReRoute()
		if to == "" {
			return
		}
		if to == RouteBack {
			r.Pop()
		} else {
			r.Push(to)
		}
		if receiver, ok := r.Active().(Receiver); ok {
			receiver.Receive(data)
		}
	}
}

// Layout static content as rigid, then layout the active route.
func (r *Router) Layout(gtx C) D {
	r.Update(gtx)
	return r.Active().Layout(gtx)
}

func (r *Router) Name() string {
	defer r.lock()()
	return r.Stack[len(r.Stack)-1]
}

func (r *Router) Active() View {
	defer r.lock()()
	return r.Routes[r.Stack[len(r.Stack)-1]]
}

func (r *Router) lock() func() {
	r.Lock()
	return func() { r.Unlock() }
}

// Route is an embedable type that implements routing.
type Route struct {
	Path string
	Data interface{}
}

// ReRoute tells the router to route to the named route.
func (r *Route) ReRoute() (string, interface{}) {
	defer func() { r.Path = "" }()
	return r.Path, r.Data
}

// To sets the route path.
func (r *Route) To(path string, data ...interface{}) {
	r.Path = path
	if len(data) > 0 {
		r.Data = data[0]
	}
}

// Back sets the route path to the special route "back".
// Tells the router to pop the view off the stack.
func (r *Route) Back() {
	r.Path = RouteBack
}
