package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/icons"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/nav"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/util"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/widget/style"

	"gioui.org/unit"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gui/views"
	"github.com/jackmordaunt/avisha-fn/notify"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func main() {
	db, err := func() (*storm.DB, error) {
		db, err := storm.Open(func() string {
			db, ok := os.LookupEnv("avisha_db")
			if !ok {
				db = filepath.Join("target", "avisha.db")
				_ = os.Mkdir(filepath.Dir(db), 0777)
			}
			return db
		}())
		if err != nil {
			return nil, err
		}
		if err := db.Init(&avisha.Lease{}); err != nil {
			return nil, err
		}
		if err := db.Init(&avisha.Site{}); err != nil {
			return nil, err
		}
		if err := db.Init(&avisha.Tenant{}); err != nil {
			return nil, err
		}
		return db, nil
	}()
	if err != nil {
		log.Fatalf("error: opening database: %v", err)
	}
	defer db.Close()
	api := avisha.App{
		DB:       db,
		Notifier: &notify.Console{},
	}
	w := app.NewWindow(app.Title("Avisha"), app.MinSize(unit.Dp(400), unit.Dp(400)))
	th := style.NewTheme(style.BootstrapPalette)
	ui := &UI{
		Window: w,
		Th:     th,
		Router: nav.Router{
			Routes: map[string]nav.View{
				views.RouteLease:      &views.LeaseList{App: &api, Th: th},
				views.RouteTenants:    &views.Tenants{App: &api, Th: th},
				views.RouteSites:      &views.Sites{App: &api, Th: th},
				views.RouteLeasePage:  &views.LeasePage{App: &api, Th: th},
				views.RouteTenantForm: &views.TenantForm{App: &api, Th: th},
				views.RouteSiteForm:   &views.SiteForm{App: &api, Th: th},
			},
			Stack: []string{views.RouteLease},
		},
		Rail: style.NavRail{
			Th:    th,
			Width: unit.Dp(80),
			Destinations: []style.Destination{
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
		if err := ui.Loop(); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// UI is the high level object that contains all global state.
// Anything that needs to integrate with the external system is allocated on
// this object.
type UI struct {
	*app.Window
	Th     *style.Theme
	Router nav.Router
	Rail   style.NavRail
	Modal  layout.Widget
}

func (ui *UI) Loop() error {
	var ops op.Ops
	for {
		switch event := (<-ui.Events()).(type) {
		case system.DestroyEvent:
			return event.Err
		case system.ClipboardEvent:
			fmt.Printf("clipboard: %v\n", event.Text)
		case *system.CommandEvent:
			if event.Type == system.CommandBack {
				ui.Router.Pop()
			}
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, event)
			ui.Layout(gtx)
			event.Frame(gtx.Ops)
		}
	}
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (ui *UI) Layout(gtx C) D {
	for _, d := range ui.Rail.Destinations {
		if d.Clicked() {
			ui.Router.Push(d.Route, nil)
		}
	}
	for ii := range ui.Rail.Destinations {
		d := &ui.Rail.Destinations[ii]
		d.Active = d.Route == ui.Router.Name()
	}
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return style.TopBar{
				Theme:  ui.Th.Primary(),
				Height: unit.Dp(50),
			}.Layout(
				gtx,
				func() string {
					if titled, ok := ui.Router.Active().(nav.Titled); ok {
						return titled.Title()
					}
					return ""
				}(),
				func() []layout.Widget {
					if contexter, ok := ui.Router.Active().(nav.Contexter); ok {
						return contexter.Context()
					}
					return nil
				}()...)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(
				gtx,
				layout.Rigid(func(gtx C) D {
					return ui.Rail.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Stack{}.Layout(
						gtx,
						layout.Expanded(func(gtx C) D {
							return util.DrawRect(
								gtx,
								color.NRGBA{R: 240, G: 240, B: 240, A: 255},
								gtx.Constraints.Max,
								unit.Dp(0))
						}),
						layout.Stacked(func(gtx C) D {
							return ui.Router.Layout(gtx)
						}),
						layout.Expanded(func(gtx C) D {
							if ui.Modal == nil {
								return D{}
							}
							return ui.Modal(gtx)
						}),
					)
				}),
			)
		}),
	)
}
