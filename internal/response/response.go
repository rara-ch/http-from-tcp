package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	Code200 StatusCode = 200
	Code400 StatusCode = 400
	Code500 StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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
	_, err := w.Write([]byte(statusLine + "\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		fieldLine := fmt.Sprintf("%s: %s \r\n", key, value)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
