package main

import (
	"net"
	"sync"
)

type Server struct {
	addr     string
	listener net.Listener
	quit     chan struct{}
	clients  sync.Map
}
