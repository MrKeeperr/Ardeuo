package models

import "time"

type AuditEvent struct {
	EventoID     int        `json:"evento_id"`
	NodeID       string     `json:"node_id"`
	Timestamp    time.Time  `json:"timestamp"`
	TipoEventoID int        `json:"tipo_evento_id"`
	DuracionRiego *int       `json:"duracion_riego"`
	VolumenAgua  *float64   `json:"volumen_agua"`
	UsuarioID    int        `json:"usuario_id"`
	Username     string     `json:"username"`
}
