package main

import (
	"fmt"
	"log"
	"os"

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
	// Ctx is shorthand for `layout.Context`.
	Ctx = layout.Context
	// Dims is shorthand for `layout.Dimensions`.
	Dims = layout.Dimensions
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
	router := &Router{
		Static: func(gtx Ctx, r *Router) Dims {
			gtx.Constraints.Max.Y = 80
			return style.TopBar{Theme: th}.Layout(
				gtx,
				r.Name(),
				func() (context []layout.Widget) {
					if contexter, ok := r.Active().(Contexter); ok {
						context = contexter.Context()
					}
					return context
				}()...)
		},
		Routes: map[string]View{
			"Lease":     &views.Lease{App: &api, Theme: th},
			"LeaseForm": &views.LeaseForm{App: &api, Theme: th},
		},
		Stack: []string{"Lease"},
	}
	go func() {
		if err := loop(w, router); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, r *Router) error {
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
