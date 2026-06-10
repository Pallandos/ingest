package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DataEnvelope permet de garder tes métadonnées (IP, Date) tout en gardant
// la donnée brute envoyée par le client intacte.
type DataEnvelope struct {
	ReceivedAt time.Time       `json:"received_at"`
	RemoteAddr string          `json:"remote_addr"`
	Payload    json.RawMessage `json:"payload"`
}

const MaxFileSize = 10 * 1024 * 1024

var (
	logger      = log.New(os.Stdout, "", log.LstdFlags)
	fileMutexes sync.Map
)

// getChannelMutex récupère ou crée un verrou unique pour un channel donné
func getChannelMutex(channel string) *sync.Mutex {
	m, _ := fileMutexes.LoadOrStore(channel, &sync.Mutex{})
	return m.(*sync.Mutex)
}

func handleData(w http.ResponseWriter, r *http.Request) {
	channel := r.PathValue("channel")
	if channel == "" {
		http.Error(w, "bad request: channel is required", http.StatusBadRequest)
		return
	}

	var payload json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad request: invalid json", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	envelope := DataEnvelope{
		ReceivedAt: time.Now().UTC(),
		RemoteAddr: r.RemoteAddr,
		Payload:    payload,
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "/data"
	}

	channelDir := filepath.Join(dataDir, channel)

	mu := getChannelMutex(channel)
	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(channelDir, 0755); err != nil {
		logger.Printf("error creating dir %s: %v", channelDir, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	filename := filepath.Join(channelDir, "data.ndjson")

	if info, err := os.Stat(filename); err == nil && info.Size() >= MaxFileSize {
		backupName := filepath.Join(channelDir, fmt.Sprintf("data_%d.ndjson", time.Now().Unix()))
		os.Rename(filename, backupName) // Renomme l'ancien fichier
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Printf("error opening file %s: %v", filename, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(envelope); err != nil {
		logger.Printf("error writing to file %s: %v", filename, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

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
	mux.HandleFunc("/health", handleHealth)

	mux.HandleFunc("POST /data/{channel}", handleData)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
