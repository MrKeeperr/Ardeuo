package repository

import (
	"context"
	"time"

	"ardeuo_backend/internal/models"
)

const adminQueryTimeout = 10 * time.Second

func ListTenants(ctx context.Context) ([]models.TenantSummary, error) {
	queryCtx, cancel := context.WithTimeout(ctx, adminQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT p.propietario_id,
			p.nombre_completo,
			p.email,
			p.telefono,
			p.activo,
			COUNT(DISTINCT u.usuario_id) AS usuarios_count,
			COUNT(DISTINCT c.cultivo_id) AS cultivos_count,
			COUNT(DISTINCT n.node_id) AS nodos_count
		 FROM propietarios p
		 LEFT JOIN usuarios u ON u.propietario_id = p.propietario_id
		 LEFT JOIN cultivos c ON c.propietario_id = p.propietario_id
		 LEFT JOIN nodos n ON n.cultivo_id = c.cultivo_id
		 GROUP BY p.propietario_id, p.nombre_completo, p.email, p.telefono, p.activo
		 ORDER BY p.propietario_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := make([]models.TenantSummary, 0)
	for rows.Next() {
		var tenant models.TenantSummary
		if err := rows.Scan(
			&tenant.PropietarioID,
			&tenant.NombreCompleto,
			&tenant.Email,
			&tenant.Telefono,
			&tenant.Activo,
			&tenant.UsuariosCount,
			&tenant.CultivosCount,
			&tenant.NodosCount,
		); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

func UpdateTenantStatus(ctx context.Context, propietarioID int, activo bool) error {
	queryCtx, cancel := context.WithTimeout(ctx, adminQueryTimeout)
	defer cancel()

	commandTag, err := DB.Exec(
		queryCtx,
		`UPDATE propietarios SET activo = $1 WHERE propietario_id = $2`,
		activo,
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

func ListTenantUsage(ctx context.Context) ([]models.TenantUsage, error) {
	queryCtx, cancel := context.WithTimeout(ctx, adminQueryTimeout)
	defer cancel()

	rows, err := DB.Query(
		queryCtx,
		`SELECT p.propietario_id,
			p.nombre_completo,
			p.activo,
			COUNT(t.*) AS telemetria_rows
		 FROM propietarios p
		 LEFT JOIN cultivos c ON c.propietario_id = p.propietario_id
		 LEFT JOIN nodos n ON n.cultivo_id = c.cultivo_id
		 LEFT JOIN telemetria t ON t.node_id = n.node_id
		 GROUP BY p.propietario_id, p.nombre_completo, p.activo
		 ORDER BY telemetria_rows DESC, p.propietario_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usage := make([]models.TenantUsage, 0)
	for rows.Next() {
		var item models.TenantUsage
		if err := rows.Scan(
			&item.PropietarioID,
			&item.NombreCompleto,
			&item.Activo,
			&item.TelemetriaRows,
		); err != nil {
			return nil, err
		}
		usage = append(usage, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usage, nil
}
