package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arcimboldo/microbot/libs/discovery"
)

func main() {

	err := discovery.CheckInBlock("http://localhost:8080", "pippo", "192.168.0.11", 1234)
	if err != nil {
		log.Panic(err)
	}
	defer discovery.CheckOut("http://localhost:8080", "pippo")
	fmt.Println("Connected. Waiting 4s before exiting")

	time.Sleep(4 * time.Second)
}
