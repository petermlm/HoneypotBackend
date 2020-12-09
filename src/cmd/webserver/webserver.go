package main

import (
	"honeypot/webserver"
	"log"
)

func main() {
	if err := webserver.ServerMain(); err != nil {
		log.Println(err)
	}
}
