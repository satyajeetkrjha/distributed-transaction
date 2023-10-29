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
	inTransaction bool
	data          map[string]string
}

func (h *webHandler) Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world"))
}

func (h *webHandler) Begin(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()

	if h.inTransaction {
		http.Error(w, "Already in transaction", http.StatusConflict)
		return
	}
	h.inTransaction = true
	w.Write([]byte("Transaction executed"))
}

func main() {
	flag.Parse()

	if *server {
		h := &webHandler{
			data: make(map[string]string),
		}

		http.Handle("/", http.HandlerFunc(h.Index))
		http.Handle("/begin", http.HandlerFunc(h.Begin))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}

}
