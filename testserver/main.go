package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})

	log.Println("Test server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
