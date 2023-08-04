package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"os"
	"tunProxy/client"
	"tunProxy/gui/resource"
	"tunProxy/log"
	"tunProxy/utils"
)

func initNodeChooser(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	proxyImg := canvas.NewImageFromResource(resource.ResourceProxyIconPng)
	proxyImg.FillMode = canvas.ImageFillOriginal
	label := widget.NewLabel("代理")
	imgAndLabelBox := container.NewHBox(proxyImg, label)
	radioBox := initProxyRadioBox(config, tunClient)
	proxyScroll = container.NewScroll(radioBox)
	proxyScroll.SetMinSize(fyne.NewSize(WIDTH-20, HEIGHT-200))
	imgLabelRadioBox := container.NewVBox(imgAndLabelBox, proxyScroll)
	return imgLabelRadioBox
}
func initClientAddrEntry(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	img := canvas.NewImageFromResource(resource.ResourceIpimgPng)
	img.FillMode = canvas.ImageFillOriginal
	label := widget.NewLabel("监听地址")
	entry := widget.NewEntry()
	entry.SetText(config.ClientAddr)
	var addrToStore string
	var button *widget.Button
	button = widget.NewButton("保存并重启", func() {
		config.ClientAddr = addrToStore
		err := config.Store(utils.CONFIG_PATH)
		if err != nil {
			log.WriteLog(log.FATAL, "保存配置失败")
		}
		entry.SetPlaceHolder(config.ClientAddr)
		button.Hide()
		os.Exit(0)
	})
	button.Hide()
	entry.SetText(config.ClientAddr)
	entry.OnChanged = func(addstr string) {
		if utils.CheckIpv4(addstr) && config.ClientAddr != addstr {
			addrToStore = addstr
			button.Show()
		} else {
			button.Hide()
		}
	}
	addrBox := container.NewHBox(img, label, entry, button)
	return addrBox
}

const DIRECT = "Direct"

func initProxyRadioBox(config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	nodes := config.Nodes
	names := make([]string, len(nodes)+1)
	names[0] = DIRECT
	i := 1
	for name, _ := range nodes {
		names[i] = name
		i++
	}
	proxyRadio := widget.NewRadio(names, func(op string) {
		go func() {
			tunClient.Vmu.Lock()
			if len(op) == 0 || op == DIRECT {
				tunClient.Mode = client.DIRECT
			} else {
				node, ok := nodes[op]
				if ok {
					tunClient.Mode = client.PROXY
					tunClient.ServerAddr = node.ServerAddr
					tunClient.Iv = []byte(node.StaticKey)
				} else {
					log.WriteLog(log.FATAL, "配置文件错误，无法加载节点")
				}
			}
			tunClient.Vmu.Unlock()
		}()
	})
	if len(names) > 1 {
		proxyRadio.SetSelected(names[1])
	} else {
		proxyRadio.SetSelected(DIRECT)
	}
	proxyBox := container.NewHBox(proxyRadio)
	return proxyBox
}
