package log

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"time"
)

const (
	maxLog         = 20
	logChanCap     = 1024
	logTextChanCap = 128
	DEBUG          = 1
	INFO           = 2
	WARNING        = 3
	FATAL          = 4
)

var lavelMap = map[string]int{"Debug": DEBUG, "Info": INFO, "Warning": WARNING, "Fatal": FATAL}
var level = INFO
var logChan chan *Log
var logTextChan chan *canvas.Text

type Log struct {
	Level int
	Text  string
}
type Queue[T any] struct {
	size int
	ch   chan T
}

func newQueue[T any](size int) *Queue[T] {
	return &Queue[T]{
		size: size,
		ch:   make(chan T, size),
	}
}
func (q *Queue[T]) put(e T) (o T) {
	select {
	case q.ch <- e:
	default:
		o = <-q.ch
		q.ch <- e
	}
	return
}
func WriteLog(level int, text string) {
	log := &Log{level, text}
	select {
	case logChan <- log:
	default:
	}
}
func InitLogger() {
	logChan = make(chan *Log, logChanCap)
	logTextChan = make(chan *canvas.Text, logTextChanCap)
	go switchLog()
}
func switchLog() {
	for {
		log := <-logChan
		if log.Level >= level {
			switch log.Level {
			case DEBUG:
				logTextChan <- canvas.NewText(log.Text, color.RGBA{169, 169, 169, 255})
			case INFO:
				logTextChan <- canvas.NewText(log.Text, color.Black)
			case WARNING:
				logTextChan <- canvas.NewText(log.Text, color.RGBA{255, 165, 4, 255})
			case FATAL:
				logTextChan <- canvas.NewText(log.Text, color.RGBA{220, 20, 60, 255})
			}
		}
	}
}
func InitLevelSelectorBox() *fyne.Container {
	selector := widget.NewSelect([]string{"Debug", "Info", "Warning", "Fatal"}, func(op string) {
		level = lavelMap[op]
	})
	selector.SetSelected("Info")
	return container.NewHBox(widget.NewLabel("日志等级"), layout.NewSpacer(), selector)
}
func InitLogScroll() *container.Scroll {
	box := container.NewVBox()
	scroll := container.NewScroll(box)
	scroll.SetMinSize(fyne.NewSize(380, 450))
	go func() {
		textQ := newQueue[*canvas.Text](maxLog)
		sepQ := newQueue[*widget.Separator](maxLog)
		for i := 0; ; i++ {
			text := <-logTextChan
			sep := widget.NewSeparator()
			text1 := textQ.put(text)
			sep1 := sepQ.put(sep)
			box.Add(text)
			box.Add(sep)
			if text1 != nil && sep1 != nil {
				time.Sleep(200 * time.Millisecond)
				box.Remove(text1)
				time.Sleep(200 * time.Millisecond)
				box.Remove(sep1)
			}
		}
	}()
	return scroll
}
