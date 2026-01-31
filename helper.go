package main

import (
	"net/http"
)



func prepareResponse(w http.ResponseWriter, r *http.Request, method string) bool{
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Preflight abfangen
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return false
	}

	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	return true
}
