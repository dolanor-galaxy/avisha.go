package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackmordaunt/avisha-fn/cmd/gio/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/theme"

	"gioui.org/unit"
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
var th = func() *theme.Theme {
	return theme.New(theme.BootstrapPalette)
}()

func main() {
	db, ok := os.LookupEnv("avisha_db")
	if !ok {
		db = "target/db.json"
	}
	api := avisha.App{
		Storage: storage.FileStorage(db).
			With(&avisha.Tenant{}).
			With(&avisha.Site{}).
			With(&avisha.Lease{}).
			Format(true).
			MustLoad(),
		Notifier: &notify.Console{},
	}
	w := app.NewWindow(app.Title("Avisha"))
	page := &nav.Page{
		Th: th,
		Router: nav.Router{
			Routes: map[string]nav.View{
				views.RouteLease:      &views.Lease{App: &api, Th: th},
				views.RouteTenants:    &views.Tenants{App: &api, Th: th},
				views.RouteSites:      &views.Sites{App: &api, Th: th},
				views.RouteLeaseForm:  &views.LeaseForm{App: &api, Th: th},
				views.RouteTenantForm: &views.TenantForm{App: &api, Th: th},
				views.RouteSiteForm:   &views.SiteForm{App: &api, Th: th},
			},
			Stack: []string{views.RouteLease},
		},
		Rail: nav.Rail{
			Width: unit.Dp(80),
			Destinations: []nav.Destination{
				{
					Label: "Leases",
					Route: views.RouteLease,
					Icon:  icons.Description,
				},
				{
					Label: "Tenants",
					Route: views.RouteTenants,
					Icon:  icons.Person,
				},
				{
					Label: "Sites",
					Route: views.RouteSites,
					Icon:  icons.Home,
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
