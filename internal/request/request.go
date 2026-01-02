package request

import (
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqRaw, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	req := string(reqRaw)

	reqLineSlice := strings.Split(req, "\r\n")[0]
	reqLine := parseRequestLine(reqLineSlice)
	if reqLine.Method == "" {
		return nil, errors.New("errors in request line")
	}

	return &Request{
		RequestLine: reqLine,
	}, nil
}

func parseRequestLine(reqLine string) RequestLine {
	splitReqLine := strings.Split(reqLine, " ")
	if len(splitReqLine) != 3 {
		return RequestLine{}
	}
	method := splitReqLine[0]
	path := splitReqLine[1]
	httpType := splitReqLine[2]

	fmt.Println(method)
	fmt.Println(path)
	fmt.Println(httpType)

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
			return RequestLine{}
		}
	}

	httpVersion := strings.TrimPrefix(httpType, "HTTP/")
	if httpVersion != "1.1" {
		return RequestLine{}
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: path,
		Method:        method,
	}
}
