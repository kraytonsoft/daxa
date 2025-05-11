package main

import (
	"daxa/runtime/handlers"
	"daxa/runtime/registry"
	"log"
	"net/http"
)

func main() {
	reg := registry.New()
	handlers.Init(reg)

	// Internal Ops API
	go func() {
		http.HandleFunc("POST /deploy/source", handlers.HandleDeploySource)
		log.Println("🔧 Daxa Internal API on :36365")
		log.Fatal(http.ListenAndServe(":36365", nil))
	}()

	// Public Function Server
	go func() {
		log.Println("🌐 Daxa Public API on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	select {} // block forever
}
