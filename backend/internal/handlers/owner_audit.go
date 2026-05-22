package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"ardeuo_backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

func OwnerListAuditEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	propietarioID, ok := getPropietarioIDFromRequest(w, r)
	if !ok {
		return
	}

	limit := parseLimitParam(r)
	items, err := repository.ListAuditEventsByPropietario(r.Context(), propietarioID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar auditoría"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func OwnerListNodoAudit(w http.ResponseWriter, r *http.Request) {
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

	limit := parseLimitParam(r)
	items, err := repository.ListNodeAuditByPropietario(r.Context(), propietarioID, nodeID, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar auditoría"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
