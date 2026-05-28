package models

import "time"

type NodoReport struct {
	NodeID         string     `json:"node_id"`
	EstadoID       int        `json:"estado_id"`
	EstadoConexion string     `json:"estado_conexion"`
	UltimaLectura  *time.Time `json:"ultima_lectura"`
}

type CultivoReport struct {
	CultivoID        int         `json:"cultivo_id"`
	AliasLote        string      `json:"alias_lote"`
	Estado           string      `json:"estado"`
	ReportFrom       time.Time   `json:"from"`
	ReportTo         time.Time   `json:"to"`
	NodosCount       int         `json:"nodos_count"`
	NodosConDatos    int         `json:"nodos_con_datos"`
	TelemetriaRows   int64       `json:"telemetria_rows"`
	AvgHumedadSuelo   *float64     `json:"avg_humedad_suelo"`
	AvgTemperaturaAmb *float64     `json:"avg_temperatura_amb"`
	AvgHumedadAmb     *float64     `json:"avg_humedad_amb"`
	UltimaLectura     *time.Time   `json:"ultima_lectura"`
	Nodos             []NodoReport `json:"nodos"`
}
