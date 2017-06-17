package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/arcimboldo/microbot/libs/discovery"
)

// courtesy of https://gist.github.com/flc/6439105
var ascii_uppercase = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ascii_lowercase = []byte("abcdefghijklmnopqrstuvwxyz")
var ascii_uppercase_len = len(ascii_uppercase)
var ascii_lowercase_len = len(ascii_lowercase)

type rot13Reader struct {
	r io.Reader
}

func rot13(b byte) byte {
	pos := bytes.IndexByte(ascii_uppercase, b)
	if pos != -1 {
		return ascii_uppercase[(pos+13)%ascii_uppercase_len]
	}
	pos = bytes.IndexByte(ascii_lowercase, b)
	if pos != -1 {
		return ascii_lowercase[(pos+13)%ascii_lowercase_len]
	}
	return b
}

func (r rot13Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = rot13(p[i])
	}
	return n, err
}

func HandleData(w http.ResponseWriter, r *http.Request) {
	r13 := rot13Reader{r.Body}
	io.Copy(w, r13)
}

func main() {
	name := flag.String("n", "rot13", "Name")
	url := flag.String("u", "http://localhost:8080", "Discovery server")
	port := flag.Int("p", 2000, "Port")

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
