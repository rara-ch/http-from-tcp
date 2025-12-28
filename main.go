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

	fmt.Printf("Reading data from %s\n", filePath)
	fmt.Println("=======================================")

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	currentLine := ""

	go func() {
		defer close(lines)
		defer f.Close()

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
