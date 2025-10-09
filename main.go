package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLStore() *URLStore {
	return &URLStore{
		urls: make(map[string]string),
	}
}

func (s *URLStore) Save(shortCode, longUrl string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[shortCode] = longUrl
}

func (s *URLStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.urls[shortCode]
	return url, exists
}

var store *URLStore

func main() {
	store = NewURLStore()
	fmt.Println("Starting a server on port 8080")
	http.HandleFunc("/", homeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL shortener is up and running")
}
