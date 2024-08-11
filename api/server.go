package api

import (
	"log"
	"net/http"
)

func NewServer() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ascii-art", AsciiArtHandler)
	http.HandleFunc("/css/", Style)
	log.Println("Starting server on http://127.0.0.1:8088")

	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
