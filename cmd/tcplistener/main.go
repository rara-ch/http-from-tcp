package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		lines := getLinesChannel(connection)

		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()
		currentLine := ""
		for {
			container := make([]byte, 8)
			_, err := f.Read(container)
			if err == io.EOF {
				lines <- currentLine
				break
			}
			if err != nil {
				log.Fatalf("error other than io.EOF retured from *os.File.Read()\n")
			}

			parts := strings.Split(string(container), "\n")

			currentLine += parts[0]

			if len(parts) == 2 {
				lines <- currentLine
				currentLine = parts[1]
			}
		}
	}()
	return lines
}
