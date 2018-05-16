package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("need URL of webserver and path to geth archives as arguments")
	}
	doMist(os.Args[1:])
}
