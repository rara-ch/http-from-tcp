package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Message string
}

// func (he *HandlerError) write(w io.Writer) {
// 	response.WriteStatusLine(w, he.StatusCode)
// 	h := response.GetDefaultHeaders(len(he.Message))
// 	response.WriteHeaders(w, h)
// 	w.Write([]byte(he.Message))
// }

type Server struct {
	closed atomic.Bool

	listener net.Listener
	handler  Handler
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

	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(500)
		body := []byte("could not parse request")
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	handlerErr := s.handler(w, req)
	if handlerErr != nil {

	}

	// headers := response.GetDefaultHeaders(len(buf.Bytes()))

	// err = response.WriteStatusLine(conn, response.Code200)
	// if err != nil {
	// 	log.Printf("error writing status line: %s", err)
	// 	conn.Close()
	// }

	// err = response.WriteHeaders(conn, headers)
	// if err != nil {
	// 	log.Printf("error writing headers: %s", err)
	// 	conn.Close()
	// }

	// conn.Write(buf.Bytes())
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}
	go s.listen()
	return s, nil
}
