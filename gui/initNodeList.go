package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"image/color"
	"tunProxy/client"
	"tunProxy/gui/resource"
	"tunProxy/log"
	"tunProxy/utils"
)

func initMyNodeTextBox(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	img := canvas.NewImageFromResource(resource.ResourceNodeIconPng)
	img.FillMode = canvas.ImageFillOriginal
	txt := canvas.NewText("我的节点", color.Black)
	addbtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if addForm.Hidden {
			addForm.Show()
			delForm.Hide()
		} else {
			addForm.Hide()
		}
	})
	delbtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if delForm.Hidden {
			delForm.Show()
			addForm.Hide()
		} else {
			delForm.Hide()
		}
	})
	box := container.NewHBox(img, txt, layout.NewSpacer(), delbtn, addbtn)
	return box
}
func initDelForm(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *widget.Form {
	nameEntry := widget.NewEntry()
	nameEntry.OnChanged = func(s string) {
		nameEntry.SetPlaceHolder("")
	}
	form := widget.NewForm(
		&widget.FormItem{Text: "名称", Widget: nameEntry},
	)
	form.Hide()
	form.CancelText = "取消"
	form.SubmitText = "删除"
	form.OnCancel = func() {
		nameEntry.SetText("")
		form.Hide()
	}
	form.OnSubmit = func() {
		name := nameEntry.Text
		if len(name) == 0 {
			nameEntry.SetText("")
			nameEntry.SetPlaceHolder("请输入节点名称")
			return
		}
		if _, ok := config.Nodes[name]; !ok {
			nameEntry.SetText("")
			nameEntry.SetPlaceHolder("节点不存在")
			return
		}
		delete(config.Nodes, name)
		err := config.Store(utils.CONFIG_PATH)
		if err != nil {
			log.WriteLog(log.FATAL, "保存配置失败")
		}
		refreshProxyScroll(proxyScroll, config, tunClient)
		refreshNodesScroll(nodesScroll, app, window, config, tunClient)
		nameEntry.SetText("")
		form.Hide()
	}
	delForm = form
	return form
}
func initAddForm(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *widget.Form {
	nameEntry := widget.NewEntry()
	addrEntry := widget.NewEntry()
	staticKeyEntry := widget.NewEntry()
	nameEntry.OnChanged = func(s string) {
		nameEntry.SetPlaceHolder("")
	}
	addrEntry.OnChanged = func(s string) {
		addrEntry.SetPlaceHolder("")
	}
	staticKeyEntry.OnChanged = func(s string) {
		staticKeyEntry.SetPlaceHolder("")
	}
	form := widget.NewForm(
		&widget.FormItem{Text: "名称", Widget: nameEntry},
		&widget.FormItem{Text: "节点地址", Widget: addrEntry},
		&widget.FormItem{Text: "秘钥", Widget: staticKeyEntry},
	)
	form.Hide()
	form.CancelText = "取消"
	form.OnCancel = func() {
		nameEntry.SetText("")
		addrEntry.SetText("")
		staticKeyEntry.SetText("")
		form.Hide()
	}

	form.SubmitText = "保存"
	nodes := config.Nodes
	form.OnSubmit = func() {
		name := nameEntry.Text
		addr := addrEntry.Text
		key := staticKeyEntry.Text
		if _, ok := nodes[name]; ok || name == "Direct" {
			nameEntry.SetText("")
			nameEntry.SetPlaceHolder("节点名重复")
		}
		if len(name) == 0 {
			nameEntry.SetText("")
			nameEntry.SetPlaceHolder("请输入节点名称")
		}
		if !utils.CheckIpv4(addr) {
			addrEntry.SetText("")
			addrEntry.SetPlaceHolder("地址格式不正确 ip:端口")
		}
		if !utils.CheckKey(key) {
			staticKeyEntry.SetText("")
			staticKeyEntry.SetPlaceHolder("密钥格式不正确")
		}
		if utils.CheckIpv4(addr) && utils.CheckKey(key) {
			config.Nodes[name] = &utils.Node{name, addr, key}
			err := config.Store(utils.CONFIG_PATH)
			if err != nil {
				log.WriteLog(log.FATAL, "保存配置失败")
			}
			nameEntry.SetText("")
			addrEntry.SetText("")
			staticKeyEntry.SetText("")
			refreshProxyScroll(proxyScroll, config, tunClient)
			refreshNodesScroll(nodesScroll, app, window, config, tunClient)
			form.Hide()
		}
	}
	addForm = form
	return form
}

func nodeToBox(node *utils.Node) *fyne.Container {
	nameImg := canvas.NewImageFromResource(resource.ResourceNamePng)
	nameImg.FillMode = canvas.ImageFillOriginal
	vpsImg := canvas.NewImageFromResource(resource.ResourceVpsPng)
	vpsImg.FillMode = canvas.ImageFillOriginal
	keyImg := canvas.NewImageFromResource(resource.ResourceKeyPng)
	keyImg.FillMode = canvas.ImageFillOriginal
	box := container.NewVBox(
		container.NewHBox(nameImg, canvas.NewText(node.Name, color.Black)),
		container.NewHBox(vpsImg, canvas.NewText(node.ServerAddr, color.Black)),
		container.NewHBox(keyImg, canvas.NewText(node.StaticKey, color.Black)),
	)
	return box
}
func initNodeBoxScroll(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *container.Scroll {
	box := initNodeBox(app, window, config, tunClient)
	scroll := container.NewScroll(box)
	scroll.SetMinSize(fyne.NewSize(WIDTH-20, HEIGHT-300))
	nodesScroll = scroll
	return scroll
}

func initNodeBox(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	flag := false
	box := container.NewVBox()
	for _, node := range config.Nodes {
		if flag {
			box.Add(widget.NewSeparator())
		}
		flag = true
		box.Add(nodeToBox(node))
	}
	return box
}
