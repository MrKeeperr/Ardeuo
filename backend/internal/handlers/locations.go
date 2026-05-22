package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ardeuo_backend/internal/repository"
)

func ListPaises(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paises, err := repository.ListPaises(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar paises"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paises)
}

func ListDepartamentos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var paisID *int
	if raw := r.URL.Query().Get("pais_id"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "pais_id invalido"})
			return
		}
		paisID = &value
	}

	departamentos, err := repository.ListDepartamentos(r.Context(), paisID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar departamentos"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(departamentos)
}

func ListMunicipios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var departamentoID *int
	if raw := r.URL.Query().Get("departamento_id"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "departamento_id invalido"})
			return
		}
		departamentoID = &value
	}

	municipios, err := repository.ListMunicipios(r.Context(), departamentoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar municipios"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(municipios)
}
