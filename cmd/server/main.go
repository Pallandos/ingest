package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type LogEntry struct {
	Level   string            `json:"level"`
	Message string            `json:"message"`
	Service string            `json:"service,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
}

type LogRecord struct {
	LogEntry
	ReceivedAt time.Time `json:"received_at"`
	RemoteAddr string    `json:"remote_addr"`
}

var logger = log.New(os.Stdout, "", 0)

func handleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry LogEntry
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&entry); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if entry.Message == "" {
		http.Error(w, "bad request: message is required", http.StatusBadRequest)
		return
	}

	record := LogRecord{
		LogEntry:   entry,
		ReceivedAt: time.Now().UTC(),
		RemoteAddr: r.RemoteAddr,
	}

	out, _ := json.Marshal(record)
	logger.Println(string(out))

	w.WriteHeader(http.StatusNoContent)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/logs", handleLog)
	mux.HandleFunc("/health", handleHealth)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
