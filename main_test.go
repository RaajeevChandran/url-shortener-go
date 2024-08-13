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
	if data.AccessCount != 0 {
		t.Errorf("Expected access count to be 0, got %d", data.AccessCount)
	}
}

func TestRedirectURLBeforeExpiry(t *testing.T) {
	url := "https://google.com"
	shortURL := createShortURL(url)

	originalURL, ok := redirectURL(shortURL)
	if !ok {
		t.Errorf("Expected to find URL for short URL %s", shortURL)
	}
	if originalURL != url {
		t.Errorf("Expected original URL %s, got %s", url, originalURL)
	}

	urlStore.RLock()
	data := urlStore.urls[shortURL]
	urlStore.RUnlock()
	if data.AccessCount != 1 {
		t.Errorf("Expected access count to be 1, got %d", data.AccessCount)
	}
}

func TestRedirectURLAfterExpiry(t *testing.T) {
	url := "https://google.com"
	shortURL := createShortURL(url)
	
	// simulate expiry
	urlStore.Lock()
	data := urlStore.urls[shortURL]
	data.Expiry = time.Now().Add(-1 * time.Hour)
	urlStore.urls[shortURL] = data
	urlStore.Unlock()

	// test after expiry
	_, ok := redirectURL(shortURL)
	if ok {
		t.Error("Expected URL to be expired, but it wasn't")
	}
}

func TestGetStats(t *testing.T) {
	url := "https://google.com"
	shortURL := createShortURL(url)

	for i := 0; i < 3; i++ {
		redirectURL(shortURL)
	}

	accessCount, ok := getStats(shortURL)
	if !ok {
		t.Errorf("Expected to find stats for short URL %s", shortURL)
	}
	if accessCount != 3 {
		t.Errorf("Expected access count to be 3, got %d", accessCount)
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
	if time.Now().After(data.Expiry) {
		t.Error("Loaded URL expiry time is in the past")
	}
	if data.AccessCount != 0 {
		t.Errorf("Expected access count to be 0, got %d", data.AccessCount)
	}
}
