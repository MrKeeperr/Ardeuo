package models

type Cultivo struct {
	ID            int      `json:"cultivo_id"`
	AliasLote     string   `json:"alias_lote"`
	PropietarioID int      `json:"propietario_id"`
	TipoCultivoID int      `json:"tipo_cultivo_id"`
	DireccionID   int      `json:"direccion_id"`
	AreaHectareas *float64 `json:"area_hectareas"`
	Estado        string   `json:"estado"`
}
