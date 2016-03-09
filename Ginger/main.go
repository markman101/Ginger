package main

import (
	"Ginger/packet"
	"Ginger/socket"
	log "github.com/donnie4w/go-logger/logger"
)

func fmain() {
	logPath := getPath()
	// init log
	log.SetConsole(true)
	log.SetRollingDaily(logPath, "client.log")
	log.SetLevel(log.DEBUG)

	protocol := &packet.TransProtocolComm{}
	reactor := &ComReactor{}
	conn, err := socket.ConServer("localhost:9009")
	if err != nil {
		log.Error("conn failed")
		return
	}

	tcpClient := socket.NewTcpSocket(conn, protocol, reactor)
	tcpClient.Start()
	data := []byte{'h', 'e', 'l', 'l', 'o'}
	tcpClient.Write(data)
	select {}
}
