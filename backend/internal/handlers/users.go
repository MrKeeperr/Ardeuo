package handlers

import (
	"encoding/json"
	"net/http"

	"ardeuo_backend/internal/repository"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.ListUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error querying usuarios"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
