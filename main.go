package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	currentLine := ""

	for {
		container := make([]byte, 8)
		_, err := file.Read(container)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error other than io.EOF retured from *os.File.Read()\n")
		}

		parts := strings.Split(string(container), "\n")

		currentLine += parts[0]

		if len(parts) == 2 {
			fmt.Printf("read: %s\n", currentLine)
			currentLine = parts[1]
		}
	}
	// Print the last line as the break stops it from being printed
	fmt.Printf("read: %s\n", currentLine)
}
