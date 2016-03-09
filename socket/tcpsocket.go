package socket

import (
	"Ginger/packet"
	"Ginger/recvbuffer"
	"errors"
	log "github.com/donnie4w/go-logger/logger"
	"io"
	"net"
	"syscall"
)

type TcpSocket struct {
	*Socket        //匿名属性
	_conn          *net.TCPConn
	_readChannel   chan *packet.Packet
	_writerChannel chan *packet.Packet
	_recvBuff      *recvbuffer.RecvBuffer
}

func ConServer(hostport string) (*net.TCPConn, error) {
	remoteAddr, err_r := net.ResolveTCPAddr("tcp4", hostport)
	if err_r != nil {
		log.Error("Conn error,ResolveTCPAddr:", hostport)
		return nil, err_r
	}
	conn, err := net.DialTCP("tcp4", nil, remoteAddr)
	if err != nil {
		log.Error("Conn error,connect to", hostport)
		return nil, err
	}

	return conn, nil
}
func NewTcpSocket(conn *net.TCPConn, protocol packet.TransProtocol, reactor EventReactor) *TcpSocket {

	sock := &TcpSocket{
		Socket: &Socket{
			_socketType:   TCPSVRCONN,
			_socketState:  ESTABLISHED,
			_tprotocol:    protocol,
			_eventReactor: reactor,
		},
		_conn:          conn,
		_readChannel:   make(chan *packet.Packet, 1000),
		_writerChannel: make(chan *packet.Packet, 1000),
		_recvBuff:      recvbuffer.NewRBuffer(packet.MAX_PACK_SIZE),
	}
	return sock
}

func (sock *TcpSocket) Start() {
	go sock.to_read()

	go sock.to_write()

	go sock.dispatch()
}
func (sock *TcpSocket) to_read() {
	for sock._socketState == ESTABLISHED {
		//接收前必须重新调整接收缓冲
		sock._recvBuff.Reset()
		//if recv too big pack which is bigger than MAX_PACK_SIZE
		if sock._recvBuff.RemainLen() <= 0 {
			log.Error("Invalid pack is too big,close the session")
			sock.closeSock()
		}
		recvLen, err := sock._conn.Read(sock._recvBuff.WritePos())
		if err != nil {
			if err == syscall.EINTR || err == syscall.EAGAIN {
				log.Error("read error--")
				continue
			} else if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				//meas read Timeout
				log.Error("read time out-", err, "neterr--", neterr)
			} else {
				log.Error("read erro-close--", err)
				sock.closeSock()
				break
			}
		}
		if recvLen > 0 {
			sock._recvBuff.WriteOffsetAdd(recvLen)
			pack, decodeLen := sock._tprotocol.Decode_pack(sock._recvBuff.ReadPos())
			if pack == nil {
				continue
			}
			sock._recvBuff.ReadOffsetAdd(decodeLen)
			//触发消息处理
			sock._readChannel <- pack
		}

	} //for

	log.Error("exit read data")
}

func (sock *TcpSocket) Write(data []byte) error {

	if len(data) <= 0 {
		return errors.New("data is error")
	}
	if ESTABLISHED != sock._socketState {
		log.Error("socket 's state is not ESTABLISHED")
		return errors.New("socket is error")
	}
	//1.encode pack
	pack := sock._tprotocol.Encode_pack(data)
	//2.write notice
	select {
	case sock._writerChannel <- pack:
		return nil
	default:
		return errors.New("WriteChannel has full")
	}
}

func (sock *TcpSocket) to_write() {
	var pack *packet.Packet
	for ESTABLISHED == sock._socketState {
		pack = <-sock._writerChannel
		if pack == nil {
			log.Error("get pack error")
			continue
		}
		data := pack.Serialize()
		n, err := sock._conn.Write(data)
		if err != nil {
			if io.ErrShortWrite != err {
				log.Error("write err close--", n)
				sock.closeSock()
				break
			}
		}
		log.Info("send data success")
	}
}

func (sock *TcpSocket) closeSock() {
	close(sock._writerChannel)
	close(sock._readChannel)
	sock._socketState = CLOSED
}

func (sock *TcpSocket) dispatch() {
	var pack *packet.Packet
	for ESTABLISHED == sock._socketState {
		pack = <-sock._readChannel
		if nil != pack {
			sock._eventReactor.OnRead(sock, pack)
		}
	}
}
