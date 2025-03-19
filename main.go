package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// / fkd -> URL {id:" " , original_url :" ", short_url:" "  , creation_date:""}
var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	fmt.Println("hash encoded to string", hash)
	fmt.Print("shortened", hash[:8])
	return hash[:8]
}
func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}
func getURL(id string) (URL, error) {
	//gets longform url from shortform
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the url shortner")
}
func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	//extract the id from request url
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
func main() {
	fmt.Println("Starting url shortner")
	// generateShortURL("https://www.github.com")
	//start the server on port 3000
	///register the handler function to handle all reqauests to the root url
	// http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/redirect/", redirectURLHandler)
	http.HandleFunc("/shorten", ShortURLHandler)
	fmt.Println("server running on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("error starting the server", err)
	}
}
