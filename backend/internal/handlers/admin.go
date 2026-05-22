package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

type UpdateTenantStatusRequest struct {
	Activo bool `json:"activo"`
}

func AdminListTenants(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tenants, err := repository.ListTenants(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar propietarios"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tenants)
}

func AdminUpdateTenantStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, err := parseIDParam(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "propietario_id inválido"})
		return
	}

	var req UpdateTenantStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	if err := repository.UpdateTenantStatus(r.Context(), propietarioID, req.Activo); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Propietario no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar propietario"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"propietario_id": propietarioID,
		"activo":         req.Activo,
	})
}

func AdminTenantUsage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usage, err := repository.ListTenantUsage(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar uso"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usage)
}

func AdminHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := repository.DB.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "degraded", "database": "down"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "database": "up"})
}

func parseIDParam(raw string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid id")
	}
	return value, nil
}
