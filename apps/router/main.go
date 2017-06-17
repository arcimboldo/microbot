package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/arcimboldo/microbot/libs/discovery"
)

var (
	flagUrl  string
	flagPort int
)

func init() {
	flag.StringVar(&flagUrl, "u", "http://localhost:8080", "Discovery server")
	flag.IntVar(&flagPort, "p", 8081, "Port to listen to.")
	flag.Parse()
}

func HandleData(w http.ResponseWriter, r *http.Request) {
	services, err := discovery.ListServices(flagUrl)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	log.Printf("Receiving call from %s%s", r.RemoteAddr, r.RequestURI)

	raw, _ := ioutil.ReadAll(r.Body)
	content := string(raw)
	data := ""

	keys := []string{}
	for k, _ := range services {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s := services[k]
		for _, re := range s.Rxs {
			log.Printf("Checking regexp %s for service %s", re, s.Name)
			if re.MatchString(content) {
				log.Printf("Service %q matching regexp %q", s.Name, re)
				resp, err := http.Post(s.URL(), "text/plain", strings.NewReader(content))
				if err != nil {
					log.Printf("Error while sending message to %q: %v", s.URL(), err)
				} else {
					raw, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Error while reading message from %q: %v", s.URL(), err)
					} else {
						data = string(raw)
						content = data
					}
				}
			}
		}
	}
	fmt.Fprintf(w, data)
	// patterns, _ := discovery.ListPatterns(flagUrl)

}

func main() {
	http.HandleFunc("/", HandleData)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil))

}
