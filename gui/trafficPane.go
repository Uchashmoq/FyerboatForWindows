package gui

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"time"
	"tunProxy/client"
	"tunProxy/utils"
)

func initTrafficPane(app *fyne.App, window *fyne.Window, config *utils.Config, tunClient *client.TunClient) *fyne.Container {
	sendLabel := widget.NewLabel("↑ 0 B/s ")
	recvLabel := widget.NewLabel("↓ 0 B/s ")
	trafficPane := container.NewHBox(sendLabel, recvLabel)
	sendTraffic := (*tunClient).SendTraffic
	recvTraffic := (*tunClient).RecvTraffic
	go freshTraffic(sendTraffic, sendLabel, "↑ %s/s ", "0 B")
	go freshTraffic(recvTraffic, recvLabel, "↓ %s/s ", "0 B")
	return trafficPane
}

func freshTraffic(traffic *client.TrafficStatistician, label *widget.Label, format, zero string) {
	for traffic.IsRunning() {
		select {
		case tf := <-traffic.Traffic:
			s := client.TrafficFormat(int64(tf))
			label.SetText(fmt.Sprintf(format, s))
		case <-time.After(2 * time.Second):
			label.SetText(fmt.Sprintf(format, zero))
		}
	}
}
