package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ardeuo_backend/internal/middleware"
	"ardeuo_backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

const (
	defaultOwnerLimit = 100
	maxOwnerLimit     = 1000
)

func getPropietarioIDFromRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return 0, false
	}

	user, err := repository.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return 0, false
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar usuario"})
		return 0, false
	}

	if user.PropietarioID == nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario sin propietario"})
		return 0, false
	}

	propietario, err := repository.GetPropietarioByID(r.Context(), *user.PropietarioID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Propietario no encontrado"})
			return 0, false
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar propietario"})
		return 0, false
	}
	if !propietario.Activo {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Propietario inactivo"})
		return 0, false
	}

	return *user.PropietarioID, true
}

func parseLimitParam(r *http.Request) int {
	value := r.URL.Query().Get("limit")
	if value == "" {
		return defaultOwnerLimit
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit <= 0 {
		return defaultOwnerLimit
	}
	if limit > maxOwnerLimit {
		return maxOwnerLimit
	}

	return limit
}
