package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := repository.DB.Query(ctx, `SELECT usuario_id, username, rol_id, fecha_creacion FROM usuarios`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error querying usuarios"})
		return
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.RoleID, &user.FechaCreacion); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "error reading usuarios"})
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error iterating usuarios"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
