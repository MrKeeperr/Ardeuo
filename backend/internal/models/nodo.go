package models

import "time"

type Nodo struct {
	NodeID          string     `json:"node_id"`
	UbicacionGeo    *string    `json:"ubicacion_geo"`
	FechaInstalacion *time.Time `json:"fecha_instalacion"`
	CultivoID       int        `json:"cultivo_id"`
	EstadoID        int        `json:"estado_id"`
}
