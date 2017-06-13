package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/arcimboldo/microbot/libs/tree"
)

var data = tree.NewTree()

func HandleGet(w http.ResponseWriter, r *http.Request) {
	path := path.Clean(r.URL.Path)
	node := data.Get(path)
	if node == nil {
		log.Printf("%s GET %s - 404 not found", r.RemoteAddr, path)
		http.NotFound(w, r)
	} else if node.Value == "" {
		log.Printf("%s GET %s - 200 listing", r.RemoteAddr, path)
		for _, child := range node.Children {
			fmt.Fprintf(w, "%s\n", child.Name)
		}
	} else {
		log.Printf("%s GET %s - 200 %s", r.RemoteAddr, path, node.Value)
		fmt.Fprintf(w, node.Value)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	val, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s POST %s=%s", r.RemoteAddr, key, val)
	data.Add(key, string(val))
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	node := data.Get(r.URL.Path)
	if node != nil {
		log.Printf("%s DELETE %s - 200 %s", r.RemoteAddr, r.URL.Path, node.Value)
		fmt.Fprintf(w, node.Value)
		data.Delete(r.URL.Path)
	} else {
		log.Printf("%s DELETE %s - 404", r.RemoteAddr, r.URL.Path)
		http.NotFound(w, r)
	}
}

func HandleData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		HandleGet(w, r)
	} else if r.Method == "POST" {
		HandlePost(w, r)
	} else if r.Method == "DELETE" {
		HandleDelete(w, r)
	}
}

func main() {
	http.HandleFunc("/", HandleData)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
