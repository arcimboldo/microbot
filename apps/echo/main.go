package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arcimboldo/microbot/libs/discovery"
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func main() {
	name := flag.String("n", "echo", "Name")
	url := flag.String("u", "http://localhost:8080", "Discovery server")
	port := flag.Int("p", 2222, "Port")

	flag.Parse()

	s := discovery.NewService(*name, "192.168.0.11", *port)
	s.AddRegexp(".*")

	err := s.CheckInBlock(*url)
	if err != nil {
		log.Panic(err)
	}
	defer s.CheckOut()

	http.HandleFunc("/", HandleData)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
