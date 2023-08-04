package client

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
	"tunProxy/log"
	"tunProxy/utils"
)

type Tun struct {
	LocalConn   net.Conn
	RemoteConn  net.Conn
	ServerAddr  string
	TargetAddr  string
	SendTraffic *TrafficStatistician
	RecvTraffic *TrafficStatistician
	key         []byte
	iv          []byte
	flag        bool
	errCh       chan error
	startTime   time.Time
	bytes       int64
}

func NewTun(localConn net.Conn, serverAddr string, iv []byte, sendTraffic, recvTraffic *TrafficStatistician) *Tun {
	return &Tun{
		LocalConn:   localConn,
		ServerAddr:  serverAddr,
		SendTraffic: sendTraffic,
		RecvTraffic: recvTraffic,
		flag:        false,
		iv:          iv,
		errCh:       make(chan error),
		bytes:       0,
	}
}
func (t *Tun) Connect() error {
	dial, err := net.Dial("tcp", t.ServerAddr)
	if err != nil {
		return err
	}
	t.RemoteConn = dial
	ch := make(chan []byte, 1)
	var err1 error
	go func() {
		b := make([]byte, 32)
		_, err1 = dial.Read(b)
		ch <- b
	}()
	select {
	case <-time.After(100 * time.Second):
		return errors.New("connect time out")
	case k := <-ch:
		t.key = k
		return err1
	}
}
func (t *Tun) StartProxy() {
	t.flag = true
	t.startTime = time.Now()
	go t.recvfromBrowserEncodeAndSend()
	go t.recvfromServerDecodeAndSendToBrowser()
	go func() {
		_ = <-t.errCh
		t.flag = false
		t.ShutDown()
		endTime := time.Now()
		log.WriteLog(log.DEBUG, fmt.Sprintf("tunnel [localhost <-> %s]closed ,communication time : %.2f s, %s transmitted ", t.TargetAddr, endTime.Sub(t.startTime).Seconds(), TrafficFormat(t.bytes)))
	}()
}
func (t *Tun) ShutDown() {
	if !t.flag {
		return
	}
	_ = t.LocalConn.Close()
	_ = t.RemoteConn.Close()
}

var timeout = 1 * time.Minute

func (t *Tun) recvfromBrowserEncodeAndSend() {
	buf := make([]byte, 1024*256)
	for t.flag {
		if err := t.LocalConn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			t.errCh <- err
			return
		}
		n, err := t.LocalConn.Read(buf)
		if err != nil {
			t.errCh <- err
			return
		}
		if len(t.TargetAddr) == 0 && IsSOCKS5CmdRequest(buf[:n]) {
			addrstr, _, _, _, _, success := parseSOCKS5Request(buf[:n])
			if success {
				log.WriteLog(log.INFO, "请求连接 "+addrstr)
				t.TargetAddr = addrstr
			}
		}
		if len(t.TargetAddr) > 0 {
			log.WriteLog(log.DEBUG, fmt.Sprintf("localhost >> %s %d bytes", t.TargetAddr, n))
		}
		message := EncodeMessage(buf[:n], t.key, t.iv)
		length := len(message)
		_, err1 := t.RemoteConn.Write(message)
		t.SendTraffic.Add(length)
		t.bytes += int64(length)
		if err1 != nil {
			t.errCh <- err
			return
		}
	}
}
func (t *Tun) recvfromServerDecodeAndSendToBrowser() {
	frameDecoder := utils.NewFrameDecoder(2, binary.BigEndian.Uint32)
	go frameDecoder.Separate()
	go func() {
		for t.flag {
			buf := make([]byte, 1024*512)
			if err := t.RemoteConn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				t.errCh <- err
				frameDecoder.Stop()
				return
			}
			n, err := t.RemoteConn.Read(buf)
			frameDecoder.In <- buf[:n]
			t.RecvTraffic.Add(n)
			t.bytes += int64(n)
			if err != nil {
				t.errCh <- err
				frameDecoder.Stop()
				return
			}
		}
	}()
	go func() {
		for t.flag {
			bytes := <-frameDecoder.Out
			message := DecodeMessage(bytes, t.key, t.iv)
			if len(t.TargetAddr) > 0 {
				log.WriteLog(log.DEBUG, fmt.Sprintf("localhost << %s %d bytes", t.TargetAddr, len(message)))
			}
			_, err1 := t.LocalConn.Write(message)
			if err1 != nil {
				t.errCh <- err1
				frameDecoder.Stop()
				return
			}
		}
	}()
}

var (
	SOCKS5_NO_AUTH   = []byte{0x05, 0x00}
	SOCKS5_SUCCESS   = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	SOCKS5_FAILED    = []byte{0x05, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	SOCKS5_UNSUPPORT = []byte{0x05, 0x07, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type DirectTun struct {
	Conn        net.Conn
	RemoteConn  net.Conn
	SendTraffic *TrafficStatistician
	RecvTraffic *TrafficStatistician
	TargetAddr  string
	flag        bool
	bytes       int64
	startTime   time.Time
}

func NewDirectTun(conn net.Conn, sendTraffic, recvTraffic *TrafficStatistician) *DirectTun {
	return &DirectTun{
		Conn:        conn,
		SendTraffic: sendTraffic,
		RecvTraffic: recvTraffic,
	}
}
func (dt *DirectTun) StartProxy() {
	dt.flag = true
	bufLocal := make([]byte, 1024*1024*2)
	go func() {
		for dt.flag {
			n, err := dt.Conn.Read(bufLocal)
			if n <= 0 || err != nil {
				dt.flag = false
				_ = dt.Conn.Close()
				break
			}
			buf := bufLocal[:n]
			if IsSOCKS5AuthRequest(buf) {
				_, _ = dt.Conn.Write(SOCKS5_NO_AUTH)
			} else if IsSOCKS5CmdRequest(buf) {
				addrstr, cmd, _, addr, port, success := parseSOCKS5Request(buf)
				if !success {
					_, _ = dt.Conn.Write(SOCKS5_FAILED)
					dt.flag = false
					_ = dt.Conn.Close()
					break
				} else if cmd != 0x01 {
					_, _ = dt.Conn.Write(SOCKS5_UNSUPPORT)
					dt.flag = false
					_ = dt.Conn.Close()
					break
				} else {
					dstAddr := &net.TCPAddr{
						IP:   addr,
						Port: int(binary.BigEndian.Uint16(port)),
					}
					dstServer, err := net.DialTCP("tcp", nil, dstAddr)
					if err != nil {
						dt.flag = false
						_ = dt.Conn.Close()
						break
					}
					dt.RemoteConn = dstServer
					_, _ = dt.Conn.Write(SOCKS5_SUCCESS)
					dt.TargetAddr = addrstr
					log.WriteLog(log.INFO, "请求连接 "+addrstr)
					dt.startTime = time.Now()
					dt.startTransmit()
					break
				}
			} else {
				dt.flag = false
				_ = dt.Conn.Close()
				break
			}
		}
	}()
}
func (dt *DirectTun) startTransmit() {
	go func() {
		bufLocalToRemote := make([]byte, 1024*1024*4)
		bufRemoteToLocal := make([]byte, 1024*1024*4)
		go func() {
			for dt.flag {
				n, err := dt.Conn.Read(bufLocalToRemote)
				if n <= 0 || err != nil {
					dt.ShutDown()
					break
				}
				if len(dt.TargetAddr) > 0 {
					log.WriteLog(log.DEBUG, fmt.Sprintf("localhost >> %s %d bytes", dt.TargetAddr, n))
				}
				dt.SendTraffic.Add(n)
				dt.bytes += int64(n)
				_, err1 := dt.RemoteConn.Write(bufLocalToRemote[:n])
				if err1 != nil {
					dt.ShutDown()
					break
				}
			}
		}()
		go func() {
			for dt.flag {
				n, err := dt.RemoteConn.Read(bufRemoteToLocal)
				if n <= 0 || err != nil {
					dt.ShutDown()
					break
				}
				if len(dt.TargetAddr) > 0 {
					log.WriteLog(log.DEBUG, fmt.Sprintf("localhost << %s %d bytes", dt.TargetAddr, n))
				}
				dt.bytes += int64(n)
				dt.RecvTraffic.Add(n)
				_, err1 := dt.Conn.Write(bufRemoteToLocal[:n])
				if err1 != nil {
					dt.ShutDown()
					break
				}
			}
		}()
	}()
}
func (dt *DirectTun) ShutDown() {
	if !dt.flag {
		return
	}
	dt.flag = false
	endTime := time.Now()
	log.WriteLog(log.DEBUG, fmt.Sprintf("tunnel [localhost <-> %s]closed ,communication time : %.2f s, %s transmitted ", dt.TargetAddr, endTime.Sub(dt.startTime).Seconds(), TrafficFormat(dt.bytes)))
	_ = dt.Conn.Close()
	_ = dt.RemoteConn.Close()
}

func parseSOCKS5Request(request []byte) (addrstr string, cmd byte, addrType byte, addr []byte, port []byte, isSuccess bool) {
	if len(request) < 7 {
		return
	}
	n := len(request)
	cmd = request[1]
	addrType = request[3]
	switch addrType {
	case 0x01: // IPv4 地址
		addr = request[4 : 4+net.IPv4len]
		port = request[n-2:]
		addrstr = fmt.Sprintf("%d.%d.%d.%d:%d", request[4], request[5], request[6], request[7], binary.BigEndian.Uint16(port))
		isSuccess = true
	case 0x03: // 域名
		port = request[n-2:]
		ipAddr, err := net.ResolveIPAddr("ip", string(request[5:n-2]))
		addrstr = fmt.Sprintf("%s:%d", string(request[5:n-2]), binary.BigEndian.Uint16(port))
		if err != nil {
			return
		}
		addr = ipAddr.IP
		isSuccess = true
	case 0x04: // IPv6 地址
		addr = request[4 : 4+net.IPv6len]
		for i := 0; i < net.IPv6len; i++ {
			if i != 0 {
				addrstr = addrstr + "."
			}
			addrstr = addrstr + fmt.Sprintf("%d", request[4+i])
		}
		addrstr = addrstr + fmt.Sprintf(":%d", binary.BigEndian.Uint16(port))
		port = request[n-2:]
		isSuccess = true
	}
	return
}
func IsSOCKS5CmdRequest(data []byte) bool {
	if len(data) < 5 || data[0] != 0x05 || data[2] != 0x00 {
		return false
	}
	cmd := data[1]
	if cmd != 0x01 && cmd != 0x03 && cmd != 0x04 {
		return false
	}
	return true
}
func IsSOCKS5AuthRequest(data []byte) bool {
	if len(data) < 3 {
		return false
	}
	if data[0] != 0x05 {
		return false
	}
	nMethods := int(data[1])
	if len(data) != 2+nMethods {
		return false
	}
	return true
}
