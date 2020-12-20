package views

import (
	"log"

	"gioui.org/layout"
	"github.com/jackmordaunt/avisha.go"
	"github.com/jackmordaunt/avisha.go/cmd/gui/nav"
	"github.com/jackmordaunt/avisha.go/cmd/gui/widget/style"
)

// SettingsPage is a page for configuring global settings.
type SettingsPage struct {
	nav.Route
	App  *avisha.App
	Th   *style.Theme
	Form SettingsForm

	scroll layout.List
}

func (s *SettingsPage) Title() string {
	return "Settings"
}

func (s *SettingsPage) Receive(_ interface{}) {
	s.Load()
}

func (s *SettingsPage) Load() {
	if settings, err := s.App.LoadSettings(); err != nil {
		log.Printf("loading settings: %v", err)
	} else {
		s.Form.Load(&settings)
	}
}

func (s *SettingsPage) Layout(gtx C) D {
	if s.Form.SubmitBtn.Clicked() {
		if settings, ok := s.Form.Submit(); ok {
			if err := s.App.SaveSettings(settings); err != nil {
				log.Printf("updating settings: %v", err)
			}
			s.Load()
		}
	}
	if s.Form.CancelBtn.Clicked() {
		s.Load()
	}
	s.scroll.Axis = layout.Vertical
	return s.scroll.Layout(gtx, 1, func(gtx C, index int) D {
		return s.Form.Layout(gtx, s.Th)
	})
}
