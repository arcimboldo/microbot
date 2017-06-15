package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/arcimboldo/microbot/libs/discovery"
)

func HandleData(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, string(content))
}

func main() {
	name := flag.String("n", "echo", "Name")
	url := flag.String("u", "http://localhost:8080", "Discovery server")
	port := flag.Int("p", 2222, "Port")

	flag.Parse()

	s := discovery.NewService(*name, "192.168.0.11", *port)
	err := s.CheckInBlock(*url)
	if err != nil {
		log.Panic(err)
	}
	defer s.CheckOut()

	http.HandleFunc("/", HandleData)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
