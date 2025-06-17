package main

import (
	"fmt"
	"go1f/pkg/server"
	"log"
)

func main() {
	fmt.Println("Server is starting")
	err := server.Run()
	if err != nil {
		log.Fatalf("Server not running: %v", err)
	}
}
