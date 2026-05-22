package models

import "time"

type UmbralAlerta struct {
	ID             int       `json:"umbral_id"`
	CultivoID      int       `json:"cultivo_id"`
	HumedadMin     float64   `json:"humedad_min"`
	HumedadMax     float64   `json:"humedad_max"`
	TemperaturaMin float64   `json:"temperatura_min"`
	TemperaturaMax float64   `json:"temperatura_max"`
	HumedadAmbMin  float64   `json:"humedad_amb_min"`
	HumedadAmbMax  float64   `json:"humedad_amb_max"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
