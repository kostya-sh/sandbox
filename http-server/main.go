package main

import (
	"log"
	"net/http"
)

var test1Body []byte = []byte("Test1!")

func main() {
	http.HandleFunc("/test1", func(w http.ResponseWriter, req *http.Request) {
		w.Write(test1Body)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
