package views

import "gioui.org/layout"

type (
	Ctx  = layout.Context
	Dims = layout.Dimensions
)

type Route = string

const (
	RouteLease     Route = "lease"
	RouteLeaseForm Route = "lease-form"
)
