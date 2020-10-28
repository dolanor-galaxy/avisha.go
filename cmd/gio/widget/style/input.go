package style

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/jackmordaunt/avisha-fn/cmd/gio/widget"
)

// The same blue used at materials.io
// Note: currently highlight color is derived from `theme.Primary`.
// var blue = color.RGBA{R: 98, G: 0, B: 238, A: 255}

// TextField implements the Material Design Text Field
// described here: https://material.io/components/text-fields
type TextField struct {
	// Editor contains the edit buffer.
	widget.Editor
	// Hoverable detects mouse hovers.
	Hoverable widget.Hoverable

	// Animation state.
	label
	border
	*anim
}

type label struct {
	TextSize float32
	Inset    layout.Inset
	Smallest layout.Dimensions
}

type border struct {
	Thickness float32
	Color     color.RGBA
}

// anim is a simple state machine for switching between animation states.
type anim struct {
	Duration  time.Duration
	began     time.Time
	progress  float32
	direction state
}

type state int8

const (
	stopped state = iota
	reverse
	forward
)

func (anim *anim) start(s state) {
	if anim.direction == stopped {
		anim.began = time.Now()
		anim.direction = s
	}
}

func (anim *anim) stop() {
	if anim.direction != stopped {
		anim.direction = stopped
	}
}

func (anim *anim) update(gtx Ctx) {
	if anim.direction != stopped {
		var (
			since = time.Since(anim.began).Milliseconds()
			total = anim.Duration.Milliseconds()
		)
		op.InvalidateOp{}.Add(gtx.Ops)
		switch anim.direction {
		case forward:
			anim.progress = float32(since) / float32(total)
			if anim.progress > 1.0 {
				anim.progress = 1.0
			}
		case reverse:
			anim.progress = 1.0 - float32(since)/float32(total)
			if anim.progress < 0.0 {
				anim.progress = 0.0
			}
		}
	}
}

func (in *TextField) Update(gtx Ctx, th *material.Theme, hint string) {
	const (
		duration = time.Millisecond * 100
	)
	if in.anim == nil {
		in.anim = &anim{Duration: duration}
	}
	if in.Editor.Focused() && in.anim.progress < 1.0 {
		in.anim.start(forward)
	}
	if (!in.Editor.Focused() && in.anim.progress == 0.0) || (in.Editor.Focused() && in.anim.progress == 1.0) {
		in.anim.stop()
	}
	if !in.Editor.Focused() && in.anim.progress > 0.0 && in.Editor.Len() == 0 {
		in.anim.start(reverse)
	}
	in.anim.update(gtx)
	var (
		// Text size transitions.
		textNormal = th.TextSize
		textSmall  = th.TextSize.Scale(0.8)
		// Border color transitions.
		borderColor        = color.RGBA{A: 107}
		borderColorHovered = color.RGBA{A: 221}
		borderColorActive  = th.Color.Primary
		// Border thickness transitions.
		borderThickness       = float32(0.5)
		borderThicknessActive = float32(2.0)
	)
	in.label.TextSize = lerp(textSmall.V, textNormal.V, 1.0-in.anim.progress)
	in.border.Thickness = borderThickness
	in.border.Color = borderColor
	if in.Hoverable.Hovered() {
		in.border.Color = borderColorHovered
	}
	if in.Editor.Focused() {
		in.border.Thickness = borderThicknessActive
		in.border.Color = borderColorActive
	}
	// Calculate the dimensions of the smallest label size and store the
	// result for use in clipping.
	// Hack: Reset min constraint to 0 to avoid min == max.
	gtx.Constraints.Min.X = 0
	macro := op.Record(gtx.Ops)
	in.label.Smallest = layout.Inset{
		Left:  unit.Dp(4),
		Right: unit.Dp(4),
	}.Layout(gtx, func(gtx Ctx) Dims {
		return material.Label(th, textSmall, hint).Layout(gtx)
	})
	macro.Stop()
	labelTopInsetNormal := float32(in.label.Smallest.Size.Y) - float32(in.label.Smallest.Size.Y/4)
	labelTopInsetActive := (labelTopInsetNormal / 2 * -1) - in.border.Thickness
	in.label.Inset = layout.Inset{
		Top:  unit.Px(lerp(labelTopInsetNormal, labelTopInsetActive, in.anim.progress)),
		Left: unit.Dp(10),
	}
}

func (in *TextField) Layout(gtx Ctx, th *material.Theme, hint string) Dims {
	in.Update(gtx, th, hint)
	defer op.Push(gtx.Ops).Pop()
	// Offset accounts for label height, which sticks above the border dimensions.
	op.Offset(f32.Pt(0, float32(in.label.Smallest.Size.Y)/2)).Add(gtx.Ops)
	in.label.Inset.Layout(
		gtx,
		func(gtx Ctx) Dims {
			return layout.Inset{
				Left:  unit.Dp(4),
				Right: unit.Dp(4),
			}.Layout(gtx, func(gtx Ctx) Dims {
				label := material.Label(th, unit.Sp(in.label.TextSize), hint)
				label.Color = in.border.Color
				return label.Layout(gtx)
			})
		})
	dims := layout.Stack{}.Layout(
		gtx,
		layout.Expanded(func(gtx Ctx) Dims {
			macro := op.Record(gtx.Ops)
			dims := widget.Border{
				Color:        in.border.Color,
				Width:        unit.Dp(in.border.Thickness),
				CornerRadius: unit.Dp(4),
			}.Layout(
				gtx,
				func(gtx Ctx) Dims {
					return Dims{Size: image.Point{
						X: gtx.Constraints.Max.X,
						Y: gtx.Constraints.Min.Y,
					}}
				},
			)
			border := macro.Stop()
			if in.Editor.Focused() || in.Editor.Len() > 0 {
				// Clip a concave shape which ignores the label area.
				clips := []clip.Rect{
					{
						Max: image.Point{
							X: gtx.Px(in.label.Inset.Left),
							Y: gtx.Constraints.Min.Y,
						},
					},
					{
						Min: image.Point{
							X: gtx.Px(in.label.Inset.Left),
							Y: in.label.Smallest.Size.Y / 2,
						},
						Max: image.Point{
							X: gtx.Px(in.label.Inset.Left) + in.label.Smallest.Size.X,
							Y: gtx.Constraints.Min.Y,
						},
					},
					{
						Min: image.Point{
							X: gtx.Px(in.label.Inset.Left) + in.label.Smallest.Size.X,
						},
						Max: image.Point{
							X: gtx.Constraints.Max.X,
							Y: gtx.Constraints.Min.Y,
						},
					},
				}
				for _, c := range clips {
					stack := op.Push(gtx.Ops)
					c.Add(gtx.Ops)
					border.Add(gtx.Ops)
					stack.Pop()
				}
			} else {
				border.Add(gtx.Ops)
			}
			return dims
		}),
		layout.Stacked(func(gtx Ctx) Dims {
			return layout.UniformInset(unit.Dp(12)).Layout(
				gtx,
				func(gtx Ctx) Dims {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return material.Editor(th, &in.Editor, "").Layout(gtx)
				},
			)
		}),
		layout.Expanded(func(gtx Ctx) Dims {
			return in.Hoverable.Layout(gtx)
		}),
	)
	return Dims{
		Size: image.Point{
			X: dims.Size.X,
			Y: dims.Size.Y + in.label.Smallest.Size.Y/2,
		},
		Baseline: dims.Baseline,
	}
}

func lerp(start, end, progress float32) float32 {
	return start + (end-start)*progress
}
