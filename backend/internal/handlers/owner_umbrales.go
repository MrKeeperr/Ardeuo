package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type UmbralRequest struct {
	HumedadMin     float64 `json:"humedad_min"`
	HumedadMax     float64 `json:"humedad_max"`
	TemperaturaMin float64 `json:"temperatura_min"`
	TemperaturaMax float64 `json:"temperatura_max"`
	HumedadAmbMin  float64 `json:"humedad_amb_min"`
	HumedadAmbMax  float64 `json:"humedad_amb_max"`
}

func OwnerUpsertUmbral(w http.ResponseWriter, r *http.Request) {
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

	var req UmbralRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	if req.HumedadMax < req.HumedadMin || req.TemperaturaMax < req.TemperaturaMin || req.HumedadAmbMax < req.HumedadAmbMin {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rangos inválidos"})
		return
	}

	umbral := models.UmbralAlerta{
		HumedadMin:     req.HumedadMin,
		HumedadMax:     req.HumedadMax,
		TemperaturaMin: req.TemperaturaMin,
		TemperaturaMax: req.TemperaturaMax,
		HumedadAmbMin:  req.HumedadAmbMin,
		HumedadAmbMax:  req.HumedadAmbMax,
	}

	result, err := repository.UpsertUmbral(r.Context(), propietarioID, cultivoID, umbral)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cultivo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al guardar umbrales"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func OwnerListUmbrales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	items, err := repository.ListUmbralesByPropietario(r.Context(), propietarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar umbrales"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
