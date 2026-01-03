package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestParseState int

const (
	initialised RequestParseState = iota
	done
)

type Request struct {
	RequestLine       RequestLine
	RequestParseState RequestParseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		RequestParseState: initialised,
	}

	for req.RequestParseState != done {
		if readToIndex >= bufferSize {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.RequestParseState = done
				break
			}

			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[:readToIndex])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func parseRequestLine(req []byte) (*RequestLine, int, error) {
	indexCRLF := bytes.Index(req, []byte(crlf))
	if indexCRLF == -1 {
		return nil, 0, nil
	}

	numBytes := len(req)

	requestLineText := string(req[:indexCRLF])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, numBytes + len(crlf), nil
}

func requestLineFromString(req string) (*RequestLine, error) {
	parts := strings.Split(req, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("error: http version, request target, and method are not seperate: %s", req)
	}

	method := parts[0]
	allowedLetters := map[string]bool{
		"A": true,
		"B": true,
		"C": true,
		"D": true,
		"E": true,
		"F": true,
		"G": true,
		"H": true,
		"I": true,
		"J": true,
		"K": true,
		"L": true,
		"M": true,
		"N": true,
		"O": true,
		"P": true,
		"Q": true,
		"R": true,
		"S": true,
		"T": true,
		"U": true,
		"V": true,
		"W": true,
		"X": true,
		"Y": true,
		"Z": true,
	}
	for _, letter := range method {
		if _, ok := allowedLetters[string(letter)]; !ok {
			return nil, fmt.Errorf("error: invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	version := parts[2]
	versionParts := strings.Split(version, "/")
	if versionParts[0] != "HTTP" || versionParts[1] != "1.1" {
		return nil, fmt.Errorf("error: invalid http version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.RequestParseState == done {
		return 0, errors.New("error: can not parse request in done state")
	}

	if r.RequestParseState == initialised {
		requestLine, numBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if numBytes == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.RequestParseState = done
		return numBytes, nil
	}

	return 0, fmt.Errorf("error: state is unkown: %v", r.RequestParseState)
}
