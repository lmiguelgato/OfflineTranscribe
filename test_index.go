package main

import (
	"embed"
	"fmt"
	"log"
)

//go:embed bundle/resources/index.html
var indexHTML embed.FS

func main() {
	fmt.Println("Testing specific index.html embed:")
	
	// Try to read index.html
	content, err := indexHTML.ReadFile("bundle/resources/index.html")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("index.html size: %d bytes\n", len(content))
	fmt.Printf("First 100 chars: %s...\n", string(content[:100]))
}