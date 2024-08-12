package main

import (
	"os"
	"testing"
	"time"
)

func TestGenerateShortURL(t *testing.T) {
	shortURL := generateShortURL(6)
	if len(shortURL) != 6 {
		t.Errorf("Expected short URL length to be 6, got %d", len(shortURL))
	}
}

func TestCreateShortURL(t *testing.T) {
	url := "https://google.com"
	shortURL := createShortURL(url)
	if shortURL == "" {
		t.Error("Expected non-empty short URL")
	}

	urlStore.RLock()
	data, ok := urlStore.urls[shortURL]
	urlStore.RUnlock()
	if !ok {
		t.Errorf("Expected URL to be stored, but it wasn't")
	}
	if data.OriginalURL != url {
		t.Errorf("Expected original URL %s, got %s", url, data.OriginalURL)
	}
	if time.Now().After(data.Expiry) {
		t.Error("URL expiry time is in the past")
	}
}

func TestRedirectURL(t *testing.T) {
	url := "https://google.com"
	shortURL := createShortURL(url)

	originalURL, ok := redirectURL(shortURL)
	if !ok {
		t.Errorf("Expected to find URL for short URL %s", shortURL)
	}
	if originalURL != url {
		t.Errorf("Expected original URL %s, got %s", url, originalURL)
	}

	// simulate expiry
	urlStore.Lock()
	data := urlStore.urls[shortURL]
	data.Expiry = time.Now().Add(-1 * time.Hour)
	urlStore.urls[shortURL] = data
	urlStore.Unlock()

	// test after expiry
	_, ok = redirectURL(shortURL)
	if ok {
		t.Error("Expected URL to be expired, but it wasn't")
	}
}

func TestSaveAndLoadURLs(t *testing.T) {
	const filename = "test_urls.json"
	defer os.Remove(filename) 

	url := "https://google.com"
	shortURL := createShortURL(url)
	saveURLsToFile(filename)

	urlStore.Lock()
	urlStore.urls = make(map[string]URLData)
	urlStore.Unlock()
	loadURLsFromFile(filename)

	urlStore.RLock()
	data, ok := urlStore.urls[shortURL]
	urlStore.RUnlock()
	if !ok {
		t.Errorf("Expected URL to be loaded from file, but it wasn't")
	}
	if data.OriginalURL != url {
		t.Errorf("Expected original URL %s, got %s", url, data.OriginalURL)
	}
}
