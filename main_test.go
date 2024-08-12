package main

import (
	"os"
	"testing"
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
	originalURL, ok := urlStore.urls[shortURL]
	urlStore.RUnlock()
	if !ok {
		t.Errorf("Expected URL to be stored, but it wasn't")
	}
	if originalURL != url {
		t.Errorf("Expected original URL %s, got %s", url, originalURL)
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
}

func TestSaveAndLoadURLs(t *testing.T) {
	const filename = "test_urls.json"
	defer os.Remove(filename) 

	url := "https://google.com"
	shortURL := createShortURL(url)
	saveURLsToFile(filename)

	urlStore.Lock()
	urlStore.urls = make(map[string]string)
	urlStore.Unlock()
	loadURLsFromFile(filename)

	urlStore.RLock()
	loadedURL, ok := urlStore.urls[shortURL]
	urlStore.RUnlock()
	if !ok {
		t.Errorf("Expected URL to be loaded from file, but it wasn't")
	}
	if loadedURL != url {
		t.Errorf("Expected original URL %s, got %s", url, loadedURL)
	}
}
