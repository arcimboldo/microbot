package main

import (
	"flag"
	"fmt"

	"github.com/arcimboldo/microbot/libs/discovery"
)

var (
	flagUrl string
)

func init() {
	flag.StringVar(&flagUrl, "u", "http://localhost:8080", "Discovery server")

	flag.Parse()
}

func main() {
	services, _ := discovery.ListServices(flagUrl)
	for _, s := range services {
		fmt.Println(s)
	}
}
