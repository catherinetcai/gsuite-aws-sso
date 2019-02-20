package http

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is a helper to generate a JSON response
func JSONResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
