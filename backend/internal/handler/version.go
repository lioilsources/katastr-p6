package handler

import (
	"encoding/json"
	"net/http"
)

const appVersion = "0.1.0"

func Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"version": appVersion,
		"backend": "go",
	})
}
