package main

import (
	"chatserver/server"
)

func main() {
	//secure := flag.Bool("secure", false, "Use TLS encryption")
	var s server.ChatServer
	s = server.NewChatServer(false)
	s.Listen("172.16.36.1:3333")
	// start the server
	s.Start()

	// 172.16.36.1
}
