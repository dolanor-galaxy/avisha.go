package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackmordaunt/avisha-fn/cmd/gio/nav"

	"gioui.org/font/gofont"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/views"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget/style"
	"github.com/jackmordaunt/avisha-fn/notify"
	"github.com/jackmordaunt/avisha-fn/storage"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

type (
	// C is shorthand for `layout.Context`.
	C = layout.Context
	// D is shorthand for `layout.Dimensions`.
	D = layout.Dimensions
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
	router := &nav.Router{
		Static: func(gtx C, r *nav.Router) D {
			gtx.Constraints.Max.Y = 80
			return style.TopBar{Theme: th}.Layout(
				gtx,
				r.Name(),
				func() (context []layout.Widget) {
					if contexter, ok := r.Active().(nav.Contexter); ok {
						context = contexter.Context()
					}
					return context
				}()...)
		},
		Routes: map[string]nav.View{
			views.RouteLease:     &views.Lease{App: &api, Theme: th},
			views.RouteLeaseForm: &views.LeaseForm{App: &api, Theme: th},
		},
		Stack: []string{views.RouteLease},
	}
	go func() {
		if err := loop(w, router); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, r *nav.Router) error {
	var ops op.Ops
	for {
		switch event := (<-w.Events()).(type) {
		case system.DestroyEvent:
			return event.Err
		case system.ClipboardEvent:
			fmt.Printf("clipboard: %v\n", event.Text)
		case *system.CommandEvent:
			if event.Type == system.CommandBack {
				r.Pop()
			}
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, event)
			r.Layout(gtx)
			event.Frame(gtx.Ops)
		}
	}
}
