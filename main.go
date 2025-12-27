package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	const filePath = "messages.txt"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("could not open %s for reading: %v\n", filePath, err)
	}
	defer file.Close()

	fmt.Printf("Reading data from %s\n", filePath)
	fmt.Println("=======================================")

	for {
		container := make([]byte, 8)
		_, err := file.Read(container)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error other than io.EOF retured from *os.File.Read()\n")
		}

		fmt.Printf("read: %s\n", string(container))
	}
}
