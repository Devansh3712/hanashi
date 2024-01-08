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

		go s.ReadMessage(conn)
		go s.Broadcast()
	}
}

func (s *Server) ReadMessage(conn net.Conn) {
	defer conn.Close()

	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read from connection: %v", err)
			break
		}

		body := strings.Trim(input, "\r\n")
		s.messages <- Message{from: conn, body: body}
	}
}

func (s *Server) Broadcast() {
	for {
		message := <-s.messages
		user, _ := s.clients.Load(message.from)
		s.clients.Range(func(key, value interface{}) bool {
			conn := key.(net.Conn)
			if conn != message.from {
				output := fmt.Sprintf(
					"%s@%v> %s\n",
					user, message.from.RemoteAddr(), message.body,
				)
				conn.Write([]byte(output))
			}
			return true
		})
	}
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		quit:     make(chan struct{}),
		messages: make(chan Message),
	}
}

func main() {
	server := NewServer(":8000")
	server.Start()
}
