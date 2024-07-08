package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var urlStore = struct {
	sync.RWMutex
	urls map[string]string
}{urls: make(map[string]string)}

func createShortURL(w http.ResponseWriter, r *http.Request) {

}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create", createShortURL).Methods("POST")
	r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")

	http.ListenAndServe(":8080", r)
}
