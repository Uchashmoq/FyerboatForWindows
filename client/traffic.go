package client

import (
	"fmt"
	"time"
)

type TrafficStatistician struct {
	in      chan int
	Traffic chan int
	ticker  *time.Ticker
	status  bool
}

func NewTrafficStatistician() *TrafficStatistician {
	ts := &TrafficStatistician{
		in:      make(chan int, 64),
		Traffic: make(chan int, 1024),
		ticker:  time.NewTicker(time.Second),
		status:  true,
	}
	go ts.run()
	return ts
}

func (ts *TrafficStatistician) IsRunning() bool {
	return ts.status
}
func (ts *TrafficStatistician) Add(size int) {
	if ts.status {
		select {
		case ts.in <- size:
		default:
		}
	}
}

func (ts *TrafficStatistician) run() {
	count := 0
	for ts.status {
		select {
		case <-ts.ticker.C:
			select {
			case ts.Traffic <- count:
				count = 0
			default:
				count = 0
			}
		default:
			count += <-ts.in
		}
	}
	ts.ticker.Stop()
}
func (ts *TrafficStatistician) Stop() {
	ts.status = false
}
func TrafficFormat(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float32(bytes)/1024.0)
	} else {
		return fmt.Sprintf("%.2f MB", float32(bytes)/1024.0/1024.0)
	}
}
