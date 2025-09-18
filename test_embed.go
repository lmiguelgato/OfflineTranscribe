package main

import (
	"embed"
	"fmt"
	"log"
)

//go:embed bundle/resources/**/*
var embeddedResources embed.FS

func main() {
	fmt.Println("Embedded files:")
	
	// List all embedded files
	err := listEmbeddedFiles(".")
	if err != nil {
		log.Fatal(err)
	}
}

func listEmbeddedFiles(dir string) error {
	entries, err := embeddedResources.ReadDir(dir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		path := dir + "/" + entry.Name()
		if dir == "." {
			path = entry.Name()
		}
		
		if entry.IsDir() {
			fmt.Printf("DIR:  %s/\n", path)
			listEmbeddedFiles(path)
		} else {
			fmt.Printf("FILE: %s\n", path)
		}
	}
	
	return nil
}