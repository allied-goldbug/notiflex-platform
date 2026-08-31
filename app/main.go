package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var counter uint64

const version = "v0.1.1"

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": version,
		"pod":     os.Getenv("HOSTNAME"),
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	id := atomic.AddUint64(&counter, 1)
	podName := os.Getenv("HOSTNAME")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":  id,
		"pod": podName,
	})
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/id", idHandler)
	http.HandleFunc("/version", versionHandler)

	port := "8080"
	fmt.Printf("notiflex-api listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
