package main

import (
	"log"
	"os"
)

// geth version
func VERSION() string {
	return "1.0"
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("need URL of webserver and path to geth archives as arguments")
	}
	doMist(os.Args[1:])
}
