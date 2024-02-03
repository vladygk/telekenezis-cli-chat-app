package main

import (
	"github.com/vladygk/chat-app/server"
)

func main() {
	server := server.Initialize()
	server.Run()
	server.StartListening(4545)
}
