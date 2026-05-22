package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"ardeuo_backend/internal/models"
	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

type OwnerCreateNodoRequest struct {
	NodeID           string     `json:"node_id"`
	UbicacionGeo     *string    `json:"ubicacion_geo"`
	FechaInstalacion *time.Time `json:"fecha_instalacion"`
	CultivoID        int        `json:"cultivo_id"`
	EstadoID         int        `json:"estado_id"`
}

type OwnerUpdateNodoRequest struct {
	UbicacionGeo     *string    `json:"ubicacion_geo"`
	FechaInstalacion *time.Time `json:"fecha_instalacion"`
	CultivoID        int        `json:"cultivo_id"`
	EstadoID         int        `json:"estado_id"`
}

func OwnerListNodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	nodos, err := repository.ListNodosByPropietario(r.Context(), propietarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar nodos"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodos)
}

func OwnerCreateNodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	var req OwnerCreateNodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	req.NodeID = strings.TrimSpace(req.NodeID)
	if req.NodeID == "" || req.CultivoID <= 0 || req.EstadoID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campos obligatorios faltantes"})
		return
	}

	nodo := models.Nodo{
		NodeID:           req.NodeID,
		UbicacionGeo:     req.UbicacionGeo,
		FechaInstalacion: req.FechaInstalacion,
		CultivoID:        req.CultivoID,
		EstadoID:         req.EstadoID,
	}

	if err := repository.CreateNodo(r.Context(), propietarioID, nodo); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cultivo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear nodo"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"node_id": req.NodeID,
		"message": "Nodo creado",
	})
}

func OwnerUpdateNodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	nodeID := strings.TrimSpace(chi.URLParam(r, "id"))
	if nodeID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "node_id inválido"})
		return
	}

	var req OwnerUpdateNodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}
	if req.CultivoID <= 0 || req.EstadoID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Campos obligatorios faltantes"})
		return
	}

	nodo := models.Nodo{
		UbicacionGeo:     req.UbicacionGeo,
		FechaInstalacion: req.FechaInstalacion,
		CultivoID:        req.CultivoID,
		EstadoID:         req.EstadoID,
	}

	if err := repository.UpdateNodo(r.Context(), propietarioID, nodeID, nodo); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Nodo no encontrado"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al actualizar nodo"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"node_id": nodeID,
		"message": "Nodo actualizado",
	})
}
