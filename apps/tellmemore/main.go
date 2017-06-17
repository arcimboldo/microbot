package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arcimboldo/microbot/libs/discovery"
)

var (
	name    string
	url     string
	port    int
	trigger string
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi! Tell me more about ")
	// Skip 'tell '
	b := make([]byte, len(trigger)+1)
	r.Body.Read(b)
	io.Copy(w, r.Body)
}

func main() {
	flag.StringVar(&name, "n", "tellmemore", "Name")
	flag.StringVar(&url, "u", "http://localhost:8080", "Discovery server")
	flag.IntVar(&port, "p", 2001, "Port")
	flag.StringVar(&trigger, "t", "tell", "Trigger")

	flag.Parse()

	s := discovery.NewService(name, "192.168.0.11", port)
	s.AddRegexp(trigger + " .*")

	err := s.CheckInBlock(url)
	if err != nil {
		log.Panic(err)
	}
	defer s.CheckOut()

	http.HandleFunc("/", HandleData)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
