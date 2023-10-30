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

	// our permanent state
	data map[string]string

	// our temporary state
	inTransaction bool
	newData       map[string]string
}

func (h *webHandler) Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world"))
}

func (h *webHandler) Begin(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()
	log.Printf("Begin")

	if h.inTransaction {
		http.Error(w, "Already in transaction", http.StatusConflict)
		return
	}
	h.inTransaction = true
	h.newData = make(map[string]string)
	w.Write([]byte("Transaction started"))
}

func (h *webHandler) Set(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()
	req.ParseForm()
	key := req.Form.Get("key")
	value := req.Form.Get("value")

	log.Printf("SET %q = %q", key, value)

	if !h.inTransaction {
		http.Error(w, "Must be in  transaction", http.StatusBadRequest)
		return
	}

	if key == " " {
		http.Error(w, "Must specify a key", http.StatusBadRequest)
		return
	}
	h.newData[key] = value
	w.Write([]byte("SET Ok"))
}

func (h *webHandler) Commit(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()
	log.Printf("Commit")

	if !h.inTransaction {
		http.Error(w, "Must be in  transaction", http.StatusBadRequest)
		return
	}

	for key, value := range h.newData {
		w.Write([]byte(key + "=" + value + "\n"))
		h.data[key] = value
	}
	h.inTransaction = false

	w.Write([]byte("Commit Successful"))
}

func (h *webHandler) RollBack(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()
	log.Printf("Rollback")

	if !h.inTransaction {
		http.Error(w, "Must be in  transaction", http.StatusBadRequest)
		return
	}

	h.inTransaction = false
	w.Write([]byte("RollBack Successful"))

}

func (h *webHandler) List(w http.ResponseWriter, req *http.Request) {
	h.Lock()
	defer h.Unlock()

	for key, value := range h.data {
		w.Write([]byte(key + "=" + value + "\n"))
	}
}

func main() {
	flag.Parse()

	if *server {
		h := &webHandler{
			data:    make(map[string]string),
			newData: make(map[string]string),
		}

		http.Handle("/", http.HandlerFunc(h.List))
		http.Handle("/set", http.HandlerFunc(h.Set))
		http.Handle("/begin", http.HandlerFunc(h.Begin))
		http.Handle("/commit", http.HandlerFunc(h.Commit))
		http.Handle("/rollback", http.HandlerFunc(h.RollBack))
		http.Handle("/list", http.HandlerFunc(h.List))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	}

}
