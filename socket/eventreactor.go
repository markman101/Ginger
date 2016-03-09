package socket

import (
	"Ginger/packet"
)

type EventReactor interface {
	OnRead(sock *TcpSocket, pack *packet.Packet)
	OnConn(sock *TcpSocket)
	OnDisconn(sock *TcpSocket)
}
