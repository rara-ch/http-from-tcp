package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, r *request.Request) *server.HandlerError {
	if r.RequestLine.RequestTarget == "/yourproblem" {
		write400Response(w)
		return nil
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		write500Response(w)
		return nil
	}

	if strings.HasPrefix(r.RequestLine.RequestTarget, "/httpbin") {
		writeChunkedResponse(w, strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin"))
		return nil
	}

	write200Response(w)
	return nil
}

func write400Response(w *response.Writer) {
	body := []byte(`
		<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>
			`)

	w.WriteStatusLine(400)

	headers := response.GetDefaultHeaders(len(body))
	headers.Override("content-type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func write500Response(w *response.Writer) {
	body := []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>
	`)

	w.WriteStatusLine(500)

	headers := response.GetDefaultHeaders(len(body))
	headers.Override("content-type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func write200Response(w *response.Writer) {
	body := []byte(`
	<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
	</html>
	`)
	w.WriteStatusLine(200)

	headers := response.GetDefaultHeaders(len(body))
	headers.Override("content-type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func writeChunkedResponse(w *response.Writer, path string) {
	httpbinRes, err := http.Get("https://httpbin.org" + path)
	if err != nil {
		write500Response(w)
	}

	err = w.WriteStatusLine(httpbinRes.StatusCode)

	httpbinRes.Header.Del("content-length")
	httpbinRes.Header.Add("transfer-encoding", "chunked")
	httpbinRes.Header.Add("Trailer", "X-Content-SHA256, X-Content-Length")
	h := headers.NewHeaders()
	for key, value := range httpbinRes.Header {
		h[key] = strings.Join(value, ", ")
	}
	w.WriteHeaders(h)

	buf := make([]byte, 1024)
	body := []byte{}

	for {
		n, err := httpbinRes.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				// TODO: Implement
			}
		}

		body = append(body, buf[:n]...)
		_, err = w.WriteChunkedBody(buf)
	}

	w.WriteChunkedBodyDone()
	t := headers.NewHeaders()
	t["X-Content-Length"] = fmt.Sprintf("%d", len(body))
	sum := sha256.Sum256(body)
	hash := sum[:]
	t["X-Content-Sha256"] = hex.EncodeToString(hash)
	w.WriteTrailers(t)
}
