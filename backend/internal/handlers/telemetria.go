package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ardeuo_backend/internal/repository"
)

func ListTelemetria(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, offset, start, end, err := parseTelemetriaQueryParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	items, err := repository.ListTelemetria(r.Context(), limit, offset, start, end)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar telemetria"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func parseTelemetriaQueryParams(r *http.Request) (*int, *int, *time.Time, *time.Time, error) {
	limit, err := parseOptionalIntParam(r, "limit", func(value int) bool { return value > 0 })
	if err != nil {
		return nil, nil, nil, nil, errors.New("limit inválido")
	}

	offset, err := parseOptionalIntParam(r, "offset", func(value int) bool { return value >= 0 })
	if err != nil {
		return nil, nil, nil, nil, errors.New("offset inválido")
	}

	start, end, err := parseDateRangeParam(r)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return limit, offset, start, end, nil
}

func parseOptionalIntParam(r *http.Request, key string, valid func(int) bool) (*int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || !valid(parsed) {
		return nil, errors.New("invalid param")
	}

	return &parsed, nil
}

func parseDateRangeParam(r *http.Request) (*time.Time, *time.Time, error) {
	fromValue := strings.TrimSpace(r.URL.Query().Get("from"))
	toValue := strings.TrimSpace(r.URL.Query().Get("to"))
	if fromValue == "" && toValue == "" {
		return nil, nil, nil
	}
	if fromValue == "" || toValue == "" {
		return nil, nil, errors.New("from y to son requeridos")
	}

	start, err := time.Parse(time.RFC3339, fromValue)
	if err != nil {
		return nil, nil, errors.New("from inválido")
	}
	end, err := time.Parse(time.RFC3339, toValue)
	if err != nil {
		return nil, nil, errors.New("to inválido")
	}
	if !end.After(start) {
		return nil, nil, errors.New("to debe ser mayor que from")
	}

	return &start, &end, nil
}
