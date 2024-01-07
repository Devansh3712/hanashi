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
			log.Printf("unable to read from connection: %v",err)
			continue
		}
		username := strings.Trim(input, "\r\n")
		s.clients.Store(username, conn)

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
		s.Broadcast(message, user)
	}
}

func (s *Server) Broadcast(message string, user string) {
	s.clients.Range(func(key, value interface{}) bool {
		username := key.(string)
		conn := value.(net.Conn)
		if username != user {
			conn.Write([]byte(fmt.Sprintf("%s> %s\n", user, message)))
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
