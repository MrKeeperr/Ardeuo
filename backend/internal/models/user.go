package models

import "time"

type User struct {
	ID            int       `json:"usuario_id"`
	Username      string    `json:"username"`
	PasswordHash  string    `json:"-"` // El guion evita que la contraseña viaje al frontend
	RoleID        int       `json:"rol_id"`
	PropietarioID *int      `json:"propietario_id"`
	Activo        bool      `json:"activo"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

type Role struct {
	ID   int    `json:"rol_id"`
	Name string `json:"nombre_rol"`
}
