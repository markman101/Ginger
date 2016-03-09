package main

import (
	"Ginger/packet"
	"Ginger/socket"
	log "github.com/donnie4w/go-logger/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ComReactor struct {
}

func (reactor *ComReactor) OnRead(sock *socket.TcpSocket, pack *packet.Packet) {
	data := pack.GetData()
	log.Info("Recv:", data)
	time.Sleep(1)
	sock.Write(data)
}

func (reactor *ComReactor) OnConn(sock *socket.TcpSocket) {

}

func (reactor *ComReactor) OnDisconn(sock *socket.TcpSocket) {

}

//get app exec path
func getPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}
func main() {
	logPath := getPath()
	// init log
	log.SetConsole(true)
	log.SetRollingDaily(logPath, "network.log")
	log.SetLevel(log.DEBUG)

	protocol := &packet.TransProtocolComm{}
	reactor := &ComReactor{}
	log.Debug("protocol addr is-", protocol)
	log.Debug("reactor addr is-", reactor)
	tcpServer := socket.NewTcpListenSock(protocol, reactor)
	log.Debug("TcpServer start!")
	tcpServer.InitTcpServer("localhost:9009")
	select {}
}
