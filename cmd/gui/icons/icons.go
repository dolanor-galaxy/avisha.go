package icons

import (
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var ArrowUp *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationArrowDropUp)
	return icon
}()

var ArrowDown *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationArrowDropDown)
	return icon
}()

var NavigationArrowBack *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationArrowBack)
	return icon
}()

var Clear *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentClear)
	return icon
}()

var Reply *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentReply)
	return icon
}()

var NavigationCancel *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationCancel)
	return icon
}()

var Send *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentSend)
	return icon
}()

var Add *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentAdd)
	return icon
}()

var Copy *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentContentCopy)
	return icon
}()

var Paste *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentContentPaste)
	return icon
}()

var Filter *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentFilterList)
	return icon
}()

var Menu *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationMenu)
	return icon
}()

var Server *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionDNS)
	return icon
}()

var Settings *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionSettings)
	return icon
}()

var Chat *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.CommunicationChat)
	return icon
}()

var Identity *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionPermIdentity)
	return icon
}()

var Home *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHome)
	return icon
}()

var Edit *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ImageEdit)
	return icon
}()

var Person *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.SocialPerson)
	return icon
}()

var Description *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionDescription)
	return icon
}()
