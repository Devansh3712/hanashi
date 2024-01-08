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
	addr     string
	listener net.Listener
	quit     chan struct{}
	clients  sync.Map
	messages chan Message
}
