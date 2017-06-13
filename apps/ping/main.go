package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/arcimboldo/microbot/libs/discovery"
)

func main() {
	name := flag.String("n", "pippo", "Name")
	flag.Parse()

	s := discovery.NewService(*name, "192.168.0.11", 1234)
	err := s.CheckInBlock("http://localhost:8080")
	if err != nil {
		log.Panic(err)
	}
	defer s.CheckOut("http://localhost:8080")
	fmt.Println("Connected. Waiting 4s before exiting")

	time.Sleep(4 * time.Second)
}
