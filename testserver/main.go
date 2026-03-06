package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request from %s", r.Method, r.RemoteAddr)

		resp := Response{
			Status:  "Operational",
			Version: "v3.0.0",
			Message: "Trademarkia Hot Reload Engine is active!",
			Time:    time.Now().Format(time.Kitchen),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("🚀 API Server booting up...")
	log.Println("📡 Listening on http://localhost:8080/api/health")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server crashed:", err)
	}
}
