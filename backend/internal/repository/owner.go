package repository

import (
	"context"
	"database/sql"
	"time"

	"ardeuo_backend/internal/models"
)

const ownerQueryTimeout = 8 * time.Second

func ListUsersByPropietario(ctx context.Context, propietarioID int) ([]models.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT usuario_id, username, rol_id, propietario_id, activo, fecha_creacion
		 FROM usuarios
		 WHERE propietario_id = $1
		 ORDER BY usuario_id`,
		propietarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.RoleID,
			&user.PropietarioID,
			&user.Activo,
			&user.FechaCreacion,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUserForPropietario(ctx context.Context, propietarioID int, username string, passwordHash string, roleID int) (int, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	var userID int
	if err := DB.QueryRow(
		queryCtx,
		`INSERT INTO usuarios (username, password_hash, rol_id, propietario_id, activo)
		 VALUES ($1, $2, $3, $4, true)
		 RETURNING usuario_id`,
		username,
		passwordHash,
		roleID,
		propietarioID,
	).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func UpdateUserRoleStatus(ctx context.Context, propietarioID int, userID int, roleID int, activo bool) error {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE usuarios
		 SET rol_id = $1, activo = $2
		 WHERE usuario_id = $3 AND propietario_id = $4`,
		roleID,
		activo,
		userID,
		propietarioID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func ListCultivosByPropietario(ctx context.Context, propietarioID int) ([]models.Cultivo, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT cultivo_id, alias_lote, propietario_id, tipo_cultivo_id, direccion_id, area_hectareas, estado
		 FROM cultivos
		 WHERE propietario_id = $1
		 ORDER BY cultivo_id`,
		propietarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cultivos := make([]models.Cultivo, 0)
	for rows.Next() {
		var cultivo models.Cultivo
		var area sql.NullFloat64
		if err := rows.Scan(
			&cultivo.ID,
			&cultivo.AliasLote,
			&cultivo.PropietarioID,
			&cultivo.TipoCultivoID,
			&cultivo.DireccionID,
			&area,
			&cultivo.Estado,
		); err != nil {
			return nil, err
		}
		if area.Valid {
			cultivo.AreaHectareas = &area.Float64
		}
		cultivos = append(cultivos, cultivo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cultivos, nil
}

func CreateCultivo(ctx context.Context, propietarioID int, cultivo models.Cultivo) (int, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	estado := cultivo.Estado
	if estado == "" {
		estado = "activo"
	}

	var cultivoID int
	if err := DB.QueryRow(
		queryCtx,
		`INSERT INTO cultivos (alias_lote, propietario_id, tipo_cultivo_id, direccion_id, area_hectareas, estado)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING cultivo_id`,
		cultivo.AliasLote,
		propietarioID,
		cultivo.TipoCultivoID,
		cultivo.DireccionID,
		cultivo.AreaHectareas,
		estado,
	).Scan(&cultivoID); err != nil {
		return 0, err
	}

	return cultivoID, nil
}

func UpdateCultivo(ctx context.Context, propietarioID int, cultivoID int, cultivo models.Cultivo) error {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	estado := cultivo.Estado
	if estado == "" {
		estado = "activo"
	}

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE cultivos
		 SET alias_lote = $1,
		     tipo_cultivo_id = $2,
		     direccion_id = $3,
		     area_hectareas = $4,
		     estado = $5
		 WHERE cultivo_id = $6 AND propietario_id = $7`,
		cultivo.AliasLote,
		cultivo.TipoCultivoID,
		cultivo.DireccionID,
		cultivo.AreaHectareas,
		estado,
		cultivoID,
		propietarioID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func CloseCultivo(ctx context.Context, propietarioID int, cultivoID int) error {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE cultivos SET estado = 'cerrado'
		 WHERE cultivo_id = $1 AND propietario_id = $2`,
		cultivoID,
		propietarioID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func ListNodosByPropietario(ctx context.Context, propietarioID int) ([]models.Nodo, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT n.node_id, n.ubicacion_geo, n.fecha_instalacion, n.cultivo_id, n.estado_id
		 FROM nodos n
		 JOIN cultivos c ON c.cultivo_id = n.cultivo_id
		 WHERE c.propietario_id = $1
		 ORDER BY n.node_id`,
		propietarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodos := make([]models.Nodo, 0)
	for rows.Next() {
		var nodo models.Nodo
		var ubicacion sql.NullString
		var fecha sql.NullTime
		if err := rows.Scan(
			&nodo.NodeID,
			&ubicacion,
			&fecha,
			&nodo.CultivoID,
			&nodo.EstadoID,
		); err != nil {
			return nil, err
		}
		if ubicacion.Valid {
			value := ubicacion.String
			nodo.UbicacionGeo = &value
		}
		if fecha.Valid {
			value := fecha.Time
			nodo.FechaInstalacion = &value
		}
		nodos = append(nodos, nodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodos, nil
}

func CreateNodo(ctx context.Context, propietarioID int, nodo models.Nodo) error {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`INSERT INTO nodos (node_id, ubicacion_geo, fecha_instalacion, cultivo_id, estado_id)
		 SELECT $1, $2, $3, c.cultivo_id, $4
		 FROM cultivos c
		 WHERE c.cultivo_id = $5 AND c.propietario_id = $6`,
		nodo.NodeID,
		nodo.UbicacionGeo,
		nodo.FechaInstalacion,
		nodo.EstadoID,
		nodo.CultivoID,
		propietarioID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func UpdateNodo(ctx context.Context, propietarioID int, nodeID string, nodo models.Nodo) error {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE nodos n
		 SET ubicacion_geo = $1,
		     fecha_instalacion = $2,
		     cultivo_id = $3,
		     estado_id = $4
		 FROM cultivos ccur, cultivos cnew
		 WHERE n.node_id = $5
		   AND ccur.cultivo_id = n.cultivo_id
		   AND ccur.propietario_id = $6
		   AND cnew.cultivo_id = $3
		   AND cnew.propietario_id = $6`,
		nodo.UbicacionGeo,
		nodo.FechaInstalacion,
		nodo.CultivoID,
		nodo.EstadoID,
		nodeID,
		propietarioID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return models.ErrNotFound
	}

	return nil
}

func UpsertUmbral(ctx context.Context, propietarioID int, cultivoID int, umbral models.UmbralAlerta) (*models.UmbralAlerta, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	row := DB.QueryRow(
		queryCtx,
		`WITH owned AS (
			SELECT cultivo_id FROM cultivos WHERE cultivo_id = $1 AND propietario_id = $2
		 )
		 INSERT INTO umbrales_alerta (
			cultivo_id,
			humedad_min,
			humedad_max,
			temperatura_min,
			temperatura_max,
			humedad_amb_min,
			humedad_amb_max
		 )
		 SELECT cultivo_id, $3, $4, $5, $6, $7, $8 FROM owned
		 ON CONFLICT (cultivo_id)
		 DO UPDATE SET
			humedad_min = EXCLUDED.humedad_min,
			humedad_max = EXCLUDED.humedad_max,
			temperatura_min = EXCLUDED.temperatura_min,
			temperatura_max = EXCLUDED.temperatura_max,
			humedad_amb_min = EXCLUDED.humedad_amb_min,
			humedad_amb_max = EXCLUDED.humedad_amb_max,
			updated_at = CURRENT_TIMESTAMP
		 RETURNING umbral_id, cultivo_id, humedad_min, humedad_max, temperatura_min, temperatura_max, humedad_amb_min, humedad_amb_max, created_at, updated_at`,
		cultivoID,
		propietarioID,
		umbral.HumedadMin,
		umbral.HumedadMax,
		umbral.TemperaturaMin,
		umbral.TemperaturaMax,
		umbral.HumedadAmbMin,
		umbral.HumedadAmbMax,
	)

	var result models.UmbralAlerta
	if err := row.Scan(
		&result.ID,
		&result.CultivoID,
		&result.HumedadMin,
		&result.HumedadMax,
		&result.TemperaturaMin,
		&result.TemperaturaMax,
		&result.HumedadAmbMin,
		&result.HumedadAmbMax,
		&result.CreatedAt,
		&result.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &result, nil
}

func ListUmbralesByPropietario(ctx context.Context, propietarioID int) ([]models.UmbralAlerta, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT u.umbral_id, u.cultivo_id, u.humedad_min, u.humedad_max,
			u.temperatura_min, u.temperatura_max, u.humedad_amb_min, u.humedad_amb_max,
			u.created_at, u.updated_at
		 FROM umbrales_alerta u
		 JOIN cultivos c ON c.cultivo_id = u.cultivo_id
		 WHERE c.propietario_id = $1
		 ORDER BY u.cultivo_id`,
		propietarioID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.UmbralAlerta, 0)
	for rows.Next() {
		var item models.UmbralAlerta
		if err := rows.Scan(
			&item.ID,
			&item.CultivoID,
			&item.HumedadMin,
			&item.HumedadMax,
			&item.TemperaturaMin,
			&item.TemperaturaMax,
			&item.HumedadAmbMin,
			&item.HumedadAmbMax,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func ListAuditEventsByPropietario(ctx context.Context, propietarioID int, limit int) ([]models.AuditEvent, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT e.evento_id, e.node_id, e.timestamp, e.tipo_evento_id, e.duracion_riego, e.volumen_agua,
			u.usuario_id, u.username
		 FROM eventos_actuadores e
		 JOIN usuarioxevento ux ON ux.evento_id = e.evento_id
		 JOIN usuarios u ON u.usuario_id = ux.usuario_id
		 JOIN nodos n ON n.node_id = e.node_id
		 JOIN cultivos c ON c.cultivo_id = n.cultivo_id
		 WHERE c.propietario_id = $1
		 ORDER BY e.timestamp DESC
		 LIMIT $2`,
		propietarioID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.AuditEvent, 0)
	for rows.Next() {
		var item models.AuditEvent
		var duracion sql.NullInt32
		var volumen sql.NullFloat64
		if err := rows.Scan(
			&item.EventoID,
			&item.NodeID,
			&item.Timestamp,
			&item.TipoEventoID,
			&duracion,
			&volumen,
			&item.UsuarioID,
			&item.Username,
		); err != nil {
			return nil, err
		}
		if duracion.Valid {
			value := int(duracion.Int32)
			item.DuracionRiego = &value
		}
		if volumen.Valid {
			value := volumen.Float64
			item.VolumenAgua = &value
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func ListNodeAuditByPropietario(ctx context.Context, propietarioID int, nodeID string, limit int) ([]models.AuditEvent, error) {
	queryCtx, cancel := context.WithTimeout(ctx, ownerQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT e.evento_id, e.node_id, e.timestamp, e.tipo_evento_id, e.duracion_riego, e.volumen_agua,
			u.usuario_id, u.username
		 FROM eventos_actuadores e
		 JOIN usuarioxevento ux ON ux.evento_id = e.evento_id
		 JOIN usuarios u ON u.usuario_id = ux.usuario_id
		 JOIN nodos n ON n.node_id = e.node_id
		 JOIN cultivos c ON c.cultivo_id = n.cultivo_id
		 WHERE c.propietario_id = $1 AND e.node_id = $2
		 ORDER BY e.timestamp DESC
		 LIMIT $3`,
		propietarioID,
		nodeID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.AuditEvent, 0)
	for rows.Next() {
		var item models.AuditEvent
		var duracion sql.NullInt32
		var volumen sql.NullFloat64
		if err := rows.Scan(
			&item.EventoID,
			&item.NodeID,
			&item.Timestamp,
			&item.TipoEventoID,
			&duracion,
			&volumen,
			&item.UsuarioID,
			&item.Username,
		); err != nil {
			return nil, err
		}
		if duracion.Valid {
			value := int(duracion.Int32)
			item.DuracionRiego = &value
		}
		if volumen.Valid {
			value := volumen.Float64
			item.VolumenAgua = &value
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
