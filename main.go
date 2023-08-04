package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	oslog "log"
	"os"
	"tunProxy/client"
	"tunProxy/gui"
	"tunProxy/log"
	"tunProxy/utils"
)

func start() {
	log.InitLogger()
	fd, err0 := os.OpenFile("./log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err0 != nil {
		log.WriteLog(log.WARNING, err0.Error())
	}
	defer fd.Close()
	os.Stdout = fd
	os.Stderr = fd
	oslog.SetOutput(fd)
	oslog.Println(gui.APP_NAME + " launch")

	if err := os.Setenv("FYNE_FONT", "msyh.ttc"); err != nil {
		log.WriteLog(log.WARNING, err.Error())
	}
	app := app.New()
	window := app.NewWindow(gui.APP_NAME)
	window.Resize(fyne.NewSize(gui.WIDTH, gui.HEIGHT))
	window.SetFixedSize(true)
	config := &utils.Config{}
	err := config.Load(utils.CONFIG_PATH)
	if err != nil {
		gui.ErrorDialog(&app, &window, fmt.Sprintf("无法加载配置 : \n%v", err.Error()))
		window.ShowAndRun()
	}
	tunClient := client.NewTunClient()
	tunClient.ClientAddr = config.ClientAddr
	tunClient.Listen()
	gui.InitWindow(&app, &window, config, tunClient)
	window.ShowAndRun()
}
func main() {
	start()
}
