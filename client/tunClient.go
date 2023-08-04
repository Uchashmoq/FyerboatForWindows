package client

import (
	"github.com/armon/go-socks5"
	"net"
	"sync"
	"tunProxy/log"
)

const (
	DIRECT = 0
	PROXY  = 1
)

type TunClient struct {
	ServerAddr   string
	ClientAddr   string
	Mode         int
	SendTraffic  *TrafficStatistician
	RecvTraffic  *TrafficStatistician
	Vmu          sync.Mutex
	directServer *socks5.Server
	listener     net.Listener
	Iv           []byte
	sigpip       chan interface{}
	status       bool
}

func NewTunClient() *TunClient {
	return &TunClient{
		Mode:        PROXY,
		SendTraffic: NewTrafficStatistician(),
		RecvTraffic: NewTrafficStatistician(),
		sigpip:      make(chan interface{}, 1),
		status:      true,
	}
}

func (t *TunClient) Listen() {
	if t.listener != nil {
		_ = t.listener.Close()
		t.listener = nil
	}
	listen, err := net.Listen("tcp", t.ClientAddr)
	if err != nil {
		log.WriteLog(log.FATAL, "监听端口失败")
	}
	t.listener = listen
}
func (t *TunClient) Accepting() {
	for {
		if !t.status {
			_ = <-t.sigpip
		}
		accept, err := t.listener.Accept()
		if err != nil {
			log.WriteLog(log.WARNING, "客户端连接异常")
		} else {
			log.WriteLog(log.DEBUG, accept.RemoteAddr().String()+" connected")
			t.Vmu.Lock()
			if len(t.ServerAddr) == 0 || len(t.ClientAddr) == 0 || t.Iv == nil || t.Mode == DIRECT {
				if t.Mode != DIRECT {
					t.Mode = DIRECT
				}
				dtun := NewDirectTun(accept, t.SendTraffic, t.RecvTraffic)
				dtun.StartProxy()
			} else {
				tun := NewTun(accept, t.ServerAddr, t.Iv, t.SendTraffic, t.RecvTraffic)
				if err := tun.Connect(); err != nil {
					log.WriteLog(log.WARNING, "无法连接到代理服务器 :"+err.Error())
					accept.Close()
					continue
				}
				tun.StartProxy()
			}
			t.Vmu.Unlock()
		}
	}
}
func (t *TunClient) Stop() {
	t.status = false
}
func (t *TunClient) Run() {
	if !t.status {
		t.status = true
		t.sigpip <- true
	}
}
