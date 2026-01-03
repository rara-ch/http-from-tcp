package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	fmt.Println("Server is running")

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("could not create tcp listener: %s", err)
	}
	defer listener.Close()

	for {

		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not create accept connection from tcp listener: %s", err)
		}
		fmt.Println("Connection has been accepted")

		req, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Printf("error occured parsing request: %s", err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		fmt.Println("connection closed")
	}
}
