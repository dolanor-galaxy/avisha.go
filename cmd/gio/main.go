package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackmordaunt/avisha-fn/cmd/gio/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"

	"gioui.org/font/gofont"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/views"
	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

// package global theme state.
var th = func() *material.Theme {
	return material.NewTheme(gofont.Collection())
}()

func main() {
	api := avisha.App{
		Storage: storage.FileStorage("target/db.json").
			With(&avisha.Tenant{}).
			With(&avisha.Site{}).
			With(&avisha.Lease{}).
			Format(true).
			MustLoad(),
		Notifier: &notify.Console{},
	}
	w := app.NewWindow(app.Title("Avisha"))
	page := &nav.Page{
		Theme: th,
		Router: nav.Router{
			Routes: map[string]nav.View{
				views.RouteLease:      &views.Lease{App: &api, Theme: th},
				views.RouteTenants:    &views.Tenants{App: &api, Theme: th},
				views.RouteLeaseForm:  &views.LeaseForm{App: &api, Theme: th},
				views.RouteTenantForm: &views.TenantForm{App: &api, Theme: th},
			},
			Stack: []string{views.RouteLease},
		},
		Rail: nav.Rail{
			Width: unit.Dp(80),
			Destinations: []nav.Destination{
				// {
				// 	Label: "Home",
				// 	Route: views.RouteLease,
				// 	Icon:  icons.Home,
				// },
				{
					Label: "Leases",
					Route: views.RouteLease,
					Icon:  icons.Edit,
				},
				{
					Label: "Tenants",
					Route: views.RouteTenants,
					Icon:  icons.Edit,
				},
				{
					Label: "Sites",
					Icon:  icons.Edit,
				},
			},
		},
	}
	go func() {
		if err := loop(w, page); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, p *nav.Page) error {
	var ops op.Ops
	for {
		switch event := (<-w.Events()).(type) {
		case system.DestroyEvent:
			return event.Err
		case system.ClipboardEvent:
			fmt.Printf("clipboard: %v\n", event.Text)
		case *system.CommandEvent:
			if event.Type == system.CommandBack {
				p.Router.Pop()
			}
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, event)
			p.Layout(gtx)
			event.Frame(gtx.Ops)
		}
	}
}
