package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/font/gofont"
	"gioui.org/widget/material"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func main() {
	// app := avisha.App{
	// 	Storage: storage.FileStorage("target/db.json").
	// 		With(&avisha.Tenant{}).
	// 		With(&avisha.Site{}).
	// 		With(&avisha.Lease{}).
	// 		Format(true).
	// 		MustLoad(),
	// 	Notifier: &notify.Console{},
	// }
	go func() {
		w := app.NewWindow(app.Title("Avisha"))
		router := &Router{
			Static: func(gtx Ctx, r *Router) Dims {
				gtx.Constraints.Max.Y = 80
				return TopBar{}.Layout(gtx, r)
			},
			Routes: map[string]Route{
				"home":  &Page{},
				"other": &Page{},
			},
			Stack: []string{"home", "other"},
		}
		if err := loop(w, router); err != nil {
			log.Fatalf("error: %v", err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type (
	// Ctx is shorthand for `layout.Context`.
	Ctx = layout.Context
	// Dims is shorthand for `layout.Dimensions`.
	Dims = layout.Dimensions
)

// package global theme state.
// Note: Can this be "easily" wrapped and passed around?
var th = material.NewTheme(gofont.Collection())

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
