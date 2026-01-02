package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqRaw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, err := parseRequestLine(reqRaw)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(req []byte) (*RequestLine, error) {
	indexCRLF := bytes.Index(req, []byte(crlf))
	if indexCRLF == -1 {
		return nil, errors.New("error: crlf does not exists in request line")
	}

	requestLineText := string(req[:indexCRLF])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(req string) (*RequestLine, error) {
	parts := strings.Split(req, " ")
	if len(parts) != 3 {
		return nil, errors.New("error: http version, request target, and method are not seperate")
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
