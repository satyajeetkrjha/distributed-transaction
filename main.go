package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var server = flag.Bool("server", false, "Run as a server")
var port = flag.Int("port", 80, "Port to listen on")

type webHandler struct {
	sync.RWMutex
	data map[string]string
}

func (h *webHandler) Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world"))
}

func main() {
	flag.Parse()

	if *server {
		h := &webHandler{
			data: make(map[string]string),
		}
		fmt.Println(h)
		http.Handle("/", http.HandlerFunc(h.Index))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}

}
