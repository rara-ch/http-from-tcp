package server

import (
	"httpfromtcp/internal/response"
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
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
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
	defer conn.Close()

	err := response.WriteStatusLine(conn, response.Code200)
	if err != nil {
		log.Printf("error writing status line: %s", err)
		conn.Close()
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("error writing headers: %s", err)
		conn.Close()
	}
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
