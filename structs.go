package main

import (
	"net"
	"sync"
)

type Message struct {
	from net.Conn
	body string
}

type Server struct {
	wg       sync.WaitGroup
	addr     string
	listener net.Listener
	clients  sync.Map
	messages chan Message
}
