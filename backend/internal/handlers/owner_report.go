package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func OwnerCultivoWeeklyReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	cultivoID, err := parseIDParam(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "cultivo_id inválido"})
		return
	}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -7)

	report, err := repository.BuildCultivoReport(r.Context(), propietarioID, cultivoID, from, now)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cultivo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar reporte"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
