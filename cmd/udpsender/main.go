package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error resolving udp address: %s", err)
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error dialing udp: %s", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error reading stdin: %s", err)
		}

		connection.Write([]byte(input))
	}

}
