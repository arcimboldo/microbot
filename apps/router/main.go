package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/arcimboldo/microbot/libs/discovery"
)

var (
	flagUrl string
)

func init() {
	flag.StringVar(&flagUrl, "u", "http://localhost:8080", "Discovery server")

	flag.Parse()
}

func HandleData(w http.ResponseWriter, r *http.Request) {
	services, err := discovery.ListServices(flagUrl)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	for _, s := range services {
		fmt.Fprintf(w, "Service %s: IP %s, port: %d\n", s.Name, s.Ip, s.Port)
	}
	// patterns, _ := discovery.ListPatterns(flagUrl)

}

func main() {
	http.HandleFunc("/", HandleData)

	log.Fatal(http.ListenAndServe(":8081", nil))

}
