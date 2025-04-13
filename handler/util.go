package handler

import (
	"encoding/json"
	"net/http"
)

func writeResponse(w http.ResponseWriter, m any) error {
	if m == nil {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(m)
}
