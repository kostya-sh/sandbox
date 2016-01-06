package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", http.FileServer(http.Dir("."))))
}
