package main

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
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
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		write500Response(w)
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

	w.WriteStatusLine(response.Code400)

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

	w.WriteStatusLine(response.Code500)

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
	w.WriteStatusLine(response.Code200)

	headers := response.GetDefaultHeaders(len(body))
	headers.Override("content-type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(body)
}
