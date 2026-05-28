package models

import "time"

type Telemetria struct {
	NodeID         string    `json:"node_id"`
	Timestamp      time.Time `json:"timestamp"`
	HumedadSuelo   float64   `json:"humedad_suelo"`
	TemperaturaAmb float64   `json:"temperatura_amb"`
	HumedadAmb     float64   `json:"humedad_amb"`
}
