package socket

import (
	"Ginger/packet"
)

//socket state
const (
	CLOSED = iota
	CONNECTING
	ESTABLISHED
	LISTEN
	SOCKSTATE_NUM
)

//socket type
const (
	TCPLISTEN = iota
	TCPSVRCONN
	TCPCLIENTCONN
	UDPSOCK
	SOCKTYPE_NUM
)

type Socket struct {
	_socketType   int
	_socketState  int
	_tprotocol    packet.TransProtocol
	_eventReactor EventReactor
}

func (sock *Socket) Reset() {
	sock._socketType = SOCKTYPE_NUM
	sock._socketState = SOCKSTATE_NUM
	sock._tprotocol = nil
	sock._eventReactor = nil
}

func (sock *Socket) SetTransProtocol(protocol packet.TransProtocol) {
	sock._tprotocol = protocol
}

func (sock *Socket) SetEventReactor(reactor EventReactor) {
	sock._eventReactor = reactor
}

func (sock *Socket) GetProtocol() packet.TransProtocol {
	return sock._tprotocol
}
