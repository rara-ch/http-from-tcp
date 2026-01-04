package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	closed atomic.Bool

	listener net.Listener
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			} else {
				log.Printf("error accepting connection: %s", err)
				continue
			}
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\nContent-Type: text/plain\nContent-Length: 13\n\nHello World!"))

	conn.Close()
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
	}
	go s.listen()
	return s, nil
}
