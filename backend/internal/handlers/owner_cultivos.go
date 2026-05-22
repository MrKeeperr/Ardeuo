package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

type OwnerCreateCultivoRequest struct {
	AliasLote     string   `json:"alias_lote"`
	TipoCultivoID int      `json:"tipo_cultivo_id"`
	DireccionID   int      `json:"direccion_id"`
	AreaHectareas *float64 `json:"area_hectareas"`
	Estado        string   `json:"estado"`
}

type OwnerUpdateCultivoRequest struct {
	AliasLote     string   `json:"alias_lote"`
	TipoCultivoID int      `json:"tipo_cultivo_id"`
	DireccionID   int      `json:"direccion_id"`
	AreaHectareas *float64 `json:"area_hectareas"`
	Estado        string   `json:"estado"`
}

func OwnerListCultivos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	cultivos, err := repository.ListCultivosByPropietario(r.Context(), propietarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar cultivos"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cultivos)
}

func OwnerCreateCultivo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	var req OwnerCreateCultivoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}
	if strings.TrimSpace(req.AliasLote) == "" || req.TipoCultivoID <= 0 || req.DireccionID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campos obligatorios faltantes"})
		return
	}

	cultivo := models.Cultivo{
		AliasLote:     strings.TrimSpace(req.AliasLote),
		TipoCultivoID: req.TipoCultivoID,
		DireccionID:   req.DireccionID,
		AreaHectareas: req.AreaHectareas,
		Estado:        strings.TrimSpace(req.Estado),
	}

	newID, err := repository.CreateCultivo(r.Context(), propietarioID, cultivo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear cultivo"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"cultivo_id": newID,
		"message":    "Cultivo creado",
	})
}

func OwnerUpdateCultivo(w http.ResponseWriter, r *http.Request) {
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

	var req OwnerUpdateCultivoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}
	if strings.TrimSpace(req.AliasLote) == "" || req.TipoCultivoID <= 0 || req.DireccionID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campos obligatorios faltantes"})
		return
	}

	cultivo := models.Cultivo{
		AliasLote:     strings.TrimSpace(req.AliasLote),
		TipoCultivoID: req.TipoCultivoID,
		DireccionID:   req.DireccionID,
		AreaHectareas: req.AreaHectareas,
		Estado:        strings.TrimSpace(req.Estado),
	}

	if err := repository.UpdateCultivo(r.Context(), propietarioID, cultivoID, cultivo); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cultivo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar cultivo"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"cultivo_id": cultivoID,
		"message":    "Cultivo actualizado",
	})
}

func OwnerCloseCultivo(w http.ResponseWriter, r *http.Request) {
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

	if err := repository.CloseCultivo(r.Context(), propietarioID, cultivoID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cultivo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al cerrar cultivo"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"cultivo_id": cultivoID,
		"message":    "Cultivo cerrado",
	})
}
