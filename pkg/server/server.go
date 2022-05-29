package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/splashercn/fireim/pkg/client"
)

type Server struct {
	clients        map[*net.Conn]*client.Client
	clientsRWMutex *sync.RWMutex
}

func NewServer() Server {
	s := Server{
		clients:        make(map[*net.Conn]*client.Client),
		clientsRWMutex: &sync.RWMutex{},
	}
	return s
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}
		go s.handleConn(&conn)
	}
}

func (s *Server) handleConn(conn *net.Conn) {
	defer func() {
		(*conn).Close()
		s.clientsRWMutex.Lock()
		delete(s.clients, conn)
		s.clientsRWMutex.Unlock()
	}()
	c := client.NewClient(conn)
	s.clientsRWMutex.Lock()
	s.clients[conn] = &c
	s.clientsRWMutex.Unlock()
	reader := bufio.NewReader(*conn)
	for {
		header := make([]byte, 2)
		n, err := io.ReadFull(reader, header)
		if err != nil || n != 2 {
			print("read header error", err)
			break
		}
		len := binary.BigEndian.Uint16(header)
		log.Print(len)
		message := make([]byte, len)
		n, err = io.ReadFull(reader, message)
		if err != nil || n != int(len) {
			print("read message error", err)
			break
		}
		log.Print(string(message))
	}
}
