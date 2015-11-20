package socket

import (
	"errors"
	"fmt"
	log "github.com/donnie4w/go-logger/logger"
	"net"
	"packet"
)

type TcpListenSocket struct {
	*Socket
	_tcpListener *net.TCPListener
	stop         chan bool //停止tcpServer
}

func NewTcpListenSock(protocol packet.TransProtocol, reactor EventReactor) *TcpListenSocket {
	tcpServer := &TcpListenSocket{
		Socket: &Socket{
			_socketType:   TCPLISTEN,
			_socketState:  LISTEN,
			_tprotocol:    protocol,
			_eventReactor: reactor,
		},
		stop: make(chan bool, 1),
	}
	return tcpServer
}
func (sock *TcpListenSocket) InitTcpServer(hostAndPort string) error {

	addr, err := net.ResolveTCPAddr("tcp4", hostAndPort)
	if err != nil {
		return errors.New("InitTcpServer erro")
	}

	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return errors.New("InitTcpServer erro")
	}

	sock._tcpListener = listener
	sock._socketType = TCPLISTEN
	sock._socketState = LISTEN
	go sock.acceptConn()

	return nil
}

func (sock *TcpListenSocket) acceptConn() error {
	if sock._socketType != TCPLISTEN || sock._socketState != LISTEN {
		return errors.New("accepet errors")
	}

	log.Debug("TcpServer begin to accept!")

	for {
		conn, err := sock._tcpListener.AcceptTCP()
		select {
		case <-sock.stop:
			return nil
		default:
			fmt.Println("nothing")
		}

		if err == nil {
			tcpSocket := NewTcpSocket(conn, sock._tprotocol, sock._eventReactor)
			tcpSocket.Start()
			log.Debug("TcpServer accept a new socket")
		}
	}

}
