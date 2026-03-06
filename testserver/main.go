package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from the hot-reloaded server! Version 2")
	})

	log.Println("Test server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
