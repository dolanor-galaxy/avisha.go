package widget

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/op"
)

// Hoverable tracks mouse hovers over some area.
type Hoverable struct {
	hovered bool
}

// Hovered if mouse has entered the area.
func (h *Hoverable) Hovered() bool {
	return h.hovered
}

// Layout Hoverable according to min constraints.
func (h *Hoverable) Layout(gtx Ctx) Dims {
	h.update(gtx)
	stack := op.Push(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   h,
		Types: pointer.Enter | pointer.Leave | pointer.Cancel,
	}.Add(gtx.Ops)
	stack.Pop()
	return Dims{Size: gtx.Constraints.Min}

}

func (h *Hoverable) update(gtx Ctx) {
	for _, event := range gtx.Events(h) {
		if event, ok := event.(pointer.Event); ok {
			switch event.Type {
			case pointer.Enter:
				h.hovered = true
			case pointer.Leave, pointer.Cancel:
				h.hovered = false
			}
		}
	}
}
