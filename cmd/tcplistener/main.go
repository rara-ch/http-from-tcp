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

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("connection closed")
	}
}
