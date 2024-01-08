package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	CMD_HELP = ":h"
	CMD_QUIT = ":q"
	CMD_WRITE = ":w"

	HELP = "Commands\n" +
	"  :w\tWrite to the server\n" + 
	"  :q\tQuit connection\n"
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
		go s.BroadcastMessage()
	}
}

func (s *Server) ReadMessage(conn net.Conn) {
	defer conn.Close()

	loop:
	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read from connection: %v", err)
			break
		}
		input = strings.Trim(input, "\r\n")

		args := strings.Split(input, " ")
		switch args[0] {
		case CMD_HELP:
			conn.Write([]byte(HELP))
		case CMD_QUIT:
			conn.Write([]byte("[+] Terminating connection\n"))
			s.clients.Delete(conn)
			break loop
		case CMD_WRITE:
			s.messages <- Message{from: conn, body: input[3:]}
		default:
			conn.Write([]byte("[-] Invalid command, use :h for list of commands\n"))
		}
	}
}

func (s *Server) BroadcastMessage() {
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
