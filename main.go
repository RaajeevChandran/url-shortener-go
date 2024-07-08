package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var urlStore = struct {
	sync.RWMutex
	urls map[string]string
}{urls: make(map[string]string)}

func generateShortURL(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL(6)
	urlStore.Lock()
	urlStore.urls[shortURL] = requestData.URL
	urlStore.Unlock()

	responseData := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: r.Host + "/" + shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	urlStore.RLock()
	originalURL, ok := urlStore.urls[shortURL]
	urlStore.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create", createShortURL).Methods("POST")
	r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")

	http.ListenAndServe(":8080", r)
}
