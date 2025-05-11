package main

import (
	"daxa/runtime/registry"
	"daxa/runtime/types"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var fnRegistry = registry.New()

func deployHandler(w http.ResponseWriter, r *http.Request) {
	var manifest types.Manifest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", 400)
		return
	}

	if err := json.Unmarshal(body, &manifest); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	for _, fn := range manifest.Functions {
		if err := fnRegistry.Register(fn); err != nil {
			log.Printf("Failed to register %s: %v", fn.ID, err)
			continue
		}

		// Register route handler
		http.HandleFunc(fn.Path, func(w http.ResponseWriter, r *http.Request) {
			if h, ok := fnRegistry.Handler(fn.Path); ok {
				h(w, r)
			} else {
				http.Error(w, "Function not active", 500)
			}
		})

		log.Printf("Registered %s → %s (version %d)", fn.ID, fn.Path, fnRegistry.Versions(fn.Path))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func rollbackHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Path    string `json:"path"`
		Version int    `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if err := fnRegistry.Rollback(data.Path, data.Version); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintf(w, "Rolled back %s to version %d", data.Path, data.Version)
}

func main() {
	http.HandleFunc("/deploy", deployHandler)
	http.HandleFunc("/rollback", rollbackHandler)

	log.Println("🚀 Daxagrid runtime listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
