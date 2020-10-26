package views

import "gioui.org/layout"

type (
	C = layout.Context
	D = layout.Dimensions
)

type Route = string

const (
	RouteLease     Route = "lease"
	RouteLeaseForm Route = "lease-form"
)
