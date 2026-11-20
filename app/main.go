
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// Note represents a note stored by our API.
type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// API stores the notes and protects them against concurrent access.
type API struct {
	mu    sync.Mutex
	notes []Note
	nextID int
}

// sendJSON sends a JSON response to the client.
func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

// healthHandler checks whether the API is running.
func (api *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	api.mu.Lock()
	count := len(api.notes)
	api.mu.Unlock()

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"notes":  count,
	})
}

// notesHandler manages note creation and listing.
func (api *API) notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.mu.Lock()
		notes := make([]Note, len(api.notes))
		copy(notes, api.notes)
		api.mu.Unlock()

		sendJSON(w, http.StatusOK, notes)

	case http.MethodPost:
		var input struct {
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if input.Content == "" {
			http.Error(w, "Content cannot be empty", http.StatusBadRequest)
			return
		}

		api.mu.Lock()
		note := Note{
			ID:      api.nextID,
			Content: input.Content,
		}
		api.nextID++
		api.notes = append(api.notes, note)
		api.mu.Unlock()

		sendJSON(w, http.StatusCreated, note)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	api := &API{
		notes:  []Note{},
		nextID: 1,
	}

	http.HandleFunc("/health", api.healthHandler)
	http.HandleFunc("/notes", api.notesHandler)

	log.Println("QuickNotes API listening on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
