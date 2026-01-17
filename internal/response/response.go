package response

import (
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"net"
	"strconv"
)

type Writer struct {
	writer io.Writer
	state  writerState
}

type writerState int

const (
	statusLineState writerState = iota
	headersState
	bodyState
	trailersState
)

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		writer: conn,
		state:  statusLineState,
	}
}

// type StatusCode int

// const (
// 	Code200 StatusCode = 200
// 	Code400 StatusCode = 400
// 	Code500 StatusCode = 500
// )

func GetStatusLine(statusCode int) string {
	statusLine := "HTTP/1.1 "
	switch statusCode {
	case 200:
		statusLine += "200 OK"
	case 400:
		statusLine += "400 Bad Request"
	case 500:
		statusLine += "500 Internal Server Error"
	default:
		statusLine += fmt.Sprintf("%v ", statusCode)
	}
	return statusLine + "\r\n"
}

func (w *Writer) WriteStatusLine(statusCode int) error {
	if w.state != statusLineState {
		return errors.New("error: wrote status line after writing headers or body")
	}
	statusLine := GetStatusLine(statusCode)
	w.writer.Write([]byte(statusLine))
	w.state = headersState
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func (w *Writer) writeHeadersLoop(h headers.Headers) error {
	for key, value := range h {
		fieldLine := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.writer.Write([]byte("\r\n"))
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != headersState {
		return errors.New("error: wrote headers before writing status line or after writing body")
	}
	err := w.writeHeadersLoop(headers)
	if err != nil {
		return err
	}
	w.state = bodyState
	return nil
}

func (w *Writer) WriteTrailers(t headers.Headers) error {
	if w.state != trailersState {
		return errors.New("error: wrote trailers before chunked body was done")
	}
	err := w.writeHeadersLoop(t)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	content := []byte(fmt.Sprintf("%X\r\n%s\r\n", len(p), p))
	return w.writer.Write(content)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.state = trailersState
	return w.writer.Write([]byte(fmt.Sprintf("%X\r\n", 0)))
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != bodyState {
		return 0, errors.New("error: wrote body before writing both status line and headers")
	}
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}
