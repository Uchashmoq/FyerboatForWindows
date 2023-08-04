package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"tunProxy/client"
	"tunProxy/utils"
)

var proxyScroll *container.Scroll

func refreshProxyScroll(pscroll *container.Scroll, config *utils.Config, tunClient *client.TunClient) {
	box := initProxyRadioBox(config, tunClient)
	pscroll.Content = box
	pscroll.Refresh()
}

var addForm *widget.Form

var delForm *widget.Form

var nodesScroll *container.Scroll

func refreshNodesScroll(ndscroll *container.Scroll, app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) {
	box := initNodeBox(app, window, config, tunClient)
	nodesScroll.Content = box
	nodesScroll.Refresh()
}
