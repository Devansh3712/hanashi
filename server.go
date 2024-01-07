package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("unable to create a tcp server: %v", err)
	}
	defer listener.Close()
	s.listener = listener

	go s.AcceptConnections()
	<-s.quit
}

func (s *Server) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %v", err)
			continue
		}

		conn.Write([]byte("Enter username: "))
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read from connection: %v", err)
			continue
		}
		username := strings.Trim(input, "\r\n")
		s.clients.Store(conn, username)

		go s.ReadMessage(conn, username)
	}
}

func (s *Server) ReadMessage(conn net.Conn, user string) {
	defer conn.Close()

	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read from connection: %v", err)
			break
		}

		message := strings.Trim(input, "\r\n")
		s.Broadcast(message, conn)
	}
}

func (s *Server) Broadcast(body string, conn net.Conn) {
	user, _ := s.clients.Load(conn)
	s.clients.Range(func(key, value interface{}) bool {
		userConn := key.(net.Conn)
		if userConn != conn {
			message := fmt.Sprintf("%s@%v> %s\n", user, conn.RemoteAddr(), body)
			userConn.Write([]byte(message))
		}
		return true
	})
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		quit: make(chan struct{}),
	}
}

func main() {
	server := NewServer(":8000")
	server.Start()
}
