package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type URLData struct {
	OriginalURL string
	Expiry      time.Time
	AccessCount int
}

var urlStore = struct {
	sync.RWMutex
	urls map[string]URLData
}{urls: make(map[string]URLData)}

const saveInterval = 5 * time.Minute 
const expiryDuration = 24 * time.Hour

func generateShortURL(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createShortURL(url string) string {
	shortURL := generateShortURL(6)
	expiry := time.Now().Add(expiryDuration)

	urlStore.Lock()
	urlStore.urls[shortURL] = URLData{
		OriginalURL: url,
		Expiry:      expiry,
		AccessCount: 0,
	}
	urlStore.Unlock()

	return shortURL
}

func redirectURL(shortURL string) (string, bool) {
	urlStore.Lock()
	defer urlStore.Unlock()
	data, ok := urlStore.urls[shortURL]

	if !ok {
		return "", false
	}

	// check if the URL has expired
	if time.Now().After(data.Expiry) {
		// URL has expired so delete it from the store
		delete(urlStore.urls, shortURL)
		return "", false
	}

	data.AccessCount++
	urlStore.urls[shortURL] = data

	return data.OriginalURL, true
}

func getStats(shortURL string) (int, bool) {
	urlStore.RLock()
	defer urlStore.RUnlock()
	data, ok := urlStore.urls[shortURL]

	if !ok {
		return 0, false
	}

	return data.AccessCount, true
}

func saveURLsToFile(filename string) {
	urlStore.RLock()
	defer urlStore.RUnlock()
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(urlStore.urls); err != nil {
		fmt.Println("Error saving URLs:", err)
	}
}

func startAutoSave(filename string) {
	for {
		time.Sleep(saveInterval)
		saveURLsToFile(filename)
	}
}

func loadURLsFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&urlStore.urls); err != nil {
		fmt.Println("Error loading URLs:", err)
	}
}

func main() {
	const filename = "urls.json"
	loadURLsFromFile(filename)

	go startAutoSave(filename)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Create Short URL")
		fmt.Println("2. Redirect to Original URL")
		fmt.Println("3. View Statistics")
		fmt.Println("4. Exit")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("Enter the URL: ")
			url, _ := reader.ReadString('\n')
			url = strings.TrimSpace(url)

			shortURL := createShortURL(url)
			fmt.Printf("Short URL: %s\n", shortURL)

		case "2":
			fmt.Print("Enter the Short URL: ")
			shortURL, _ := reader.ReadString('\n')
			shortURL = strings.TrimSpace(shortURL)

			if originalURL, ok := redirectURL(shortURL); ok {
				fmt.Printf("Original URL: %s\n", originalURL)
			} else {
				fmt.Println("Short URL not found or has expired.")
			}

		case "3":
			fmt.Print("Enter the Short URL to view statistics: ")
			shortURL, _ := reader.ReadString('\n')
			shortURL = strings.TrimSpace(shortURL)

			if accessCount, ok := getStats(shortURL); ok {
				fmt.Printf("Short URL %s has been accessed %d times.\n", shortURL, accessCount)
			} else {
				fmt.Println("Short URL not found or has expired.")
			}

		case "4":
			fmt.Println("Exiting...")
			saveURLsToFile(filename) 
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
