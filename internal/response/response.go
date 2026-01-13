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
)

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		writer: conn,
		state:  statusLineState,
	}
}

type StatusCode int

const (
	Code200 StatusCode = 200
	Code400 StatusCode = 400
	Code500 StatusCode = 500
)

func GetStatusLine(statusCode StatusCode) string {
	statusLine := "HTTP/1.1 "
	switch statusCode {
	case Code200:
		statusLine += "200 OK"
	case Code400:
		statusLine += "400 Bad Request"
	case Code500:
		statusLine += "500 Internal Server Error"
	default:
		statusLine += fmt.Sprintf("%v ", statusCode)
	}
	return statusLine + "\r\n"
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != headersState {
		return errors.New("error: wrote headers before writing status line or after writing body")
	}
	for key, value := range headers {
		fieldLine := fmt.Sprintf("%s: %s \r\n", key, value)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.writer.Write([]byte("\r\n"))
	w.state = bodyState
	return nil
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
