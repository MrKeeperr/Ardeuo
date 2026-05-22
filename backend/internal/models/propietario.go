package models

type Propietario struct {
	ID             int    `json:"propietario_id"`
	NombreCompleto string `json:"nombre_completo"`
	Email          string `json:"email"`
	Telefono       string `json:"telefono"`
	DireccionID    *int   `json:"direccion_id"`
	Activo         bool   `json:"activo"`
}

type TenantSummary struct {
	PropietarioID int    `json:"propietario_id"`
	NombreCompleto string `json:"nombre_completo"`
	Email          string `json:"email"`
	Telefono       string `json:"telefono"`
	Activo         bool   `json:"activo"`
	UsuariosCount  int64  `json:"usuarios_count"`
	CultivosCount  int64  `json:"cultivos_count"`
	NodosCount     int64  `json:"nodos_count"`
}

type TenantUsage struct {
	PropietarioID int    `json:"propietario_id"`
	NombreCompleto string `json:"nombre_completo"`
	Activo         bool   `json:"activo"`
	TelemetriaRows int64  `json:"telemetria_rows"`
}
