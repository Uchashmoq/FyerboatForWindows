package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"runtime"
	"strings"
	"tunProxy/client"
	"tunProxy/gui/resource"
	"tunProxy/log"
	"tunProxy/utils"
)

const (
	WIDTH    = 450
	HEIGHT   = 600
	APP_NAME = "Fyerboat"
)

func initTabContainer(app *fyne.App, window *fyne.Window, config *utils.Config, client *client.TunClient) *widget.TabContainer {
	clientAddrEntryBox := initClientAddrEntry(app, window, config, client)
	chooser := initNodeChooser(app, window, config, client)

	myNodeBar := initMyNodeTextBox(app, window, config, client)
	aform := initAddForm(app, window, config, client)
	dform := initDelForm(app, window, config, client)
	ndscroll := initNodeBoxScroll(app, window, config, client)

	levelSelector := log.InitLevelSelectorBox()
	logScroll := log.InitLogScroll()
	tabs := widget.NewTabContainer(
		container.NewTabItem("主页", container.NewVBox(clientAddrEntryBox, widget.NewSeparator(), chooser)),
		container.NewTabItem("节点", container.NewVBox(myNodeBar, aform, dform, widget.NewSeparator(), ndscroll)),
		container.NewTabItem("日志", container.NewVBox(levelSelector, widget.NewSeparator(), logScroll)),
	)
	return tabs
}

func InitWindow(app *fyne.App, window *fyne.Window, config *utils.Config, client *client.TunClient) {
	(*window).SetIcon(resource.ResourceSymbolPng)
	(*window).CenterOnScreen()
	symbolBox := initTop(app, window, config, client)
	tabContainer := initTabContainer(app, window, config, client)

	contain := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		symbolBox,
		tabContainer,
	)
	(*window).SetContent(contain)
	go client.Accepting()
}

func initTop(app *fyne.App, window *fyne.Window, config *utils.Config, client *client.TunClient) *fyne.Container {
	symbol := canvas.NewImageFromResource(resource.ResourceSymbolPng)
	symbol.FillMode = canvas.ImageFillOriginal
	text := canvas.NewText(fmt.Sprintf("%s for %s", APP_NAME, strings.Title(runtime.GOOS)), color.Black)
	text.TextSize = 16
	trafficPane := initTrafficPane(app, window, config, client)
	topBox := container.NewHBox(symbol, text, layout.NewSpacer(), trafficPane)
	return topBox
}
