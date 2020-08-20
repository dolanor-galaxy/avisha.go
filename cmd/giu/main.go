package main

import (
	"github.com/AllenDang/giu"
)

func main() {
	var (
		content giu.Widget = giu.Label("DashBoard")
	)

	w := giu.NewMasterWindow("Avisha", 1024, 768, 0, nil)

	w.Main(func() {
		giu.SingleWindow("Avisha", giu.Layout{
			giu.Layout{
				giu.SplitLayout("Main", giu.DirectionHorizontal, true, 200,
					giu.Layout{
						giu.ButtonV("DashBoard", -1, 40, func() {
							content = giu.Label("DashBoard")
						}),
						giu.ButtonV("Lease", -1, 40, func() {
							content = giu.Label("Lease")
						}),
						giu.ButtonV("Tenant", -1, 40, func() {
							content = giu.Label("Tenant")
						}),
						giu.ButtonV("Site", -1, 40, func() {
							content = giu.Label("Site")
						}),
					},
					content,
				),
			},
		})
	})
}
