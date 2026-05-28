package repository

import (
	"context"
	"database/sql"
	"time"

	"ardeuo_backend/internal/models"
)

const cultivoReportTimeout = 12 * time.Second

func BuildCultivoReport(ctx context.Context, propietarioID int, cultivoID int, from time.Time, to time.Time) (*models.CultivoReport, error) {
	queryCtx, cancel := context.WithTimeout(ctx, cultivoReportTimeout)
	defer cancel()

	var report models.CultivoReport
	row := DB.QueryRow(
		queryCtx,
		`SELECT c.cultivo_id, c.alias_lote, c.estado
		 FROM cultivos c
		 WHERE c.cultivo_id = $1 AND c.propietario_id = $2`,
		cultivoID,
		propietarioID,
	)
	if err := row.Scan(&report.CultivoID, &report.AliasLote, &report.Estado); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	report.ReportFrom = from
	report.ReportTo = to

	nodoRows, err := DB.Query(
		queryCtx,
		`SELECT n.node_id, n.estado_id, e.estado_dela_conexion,
			MAX(t.timestamp) AS ultima_lectura
		 FROM nodos n
		 JOIN estados_conexion e ON e.estado_id = n.estado_id
		 LEFT JOIN telemetria t
			ON t.node_id = n.node_id
			AND t.timestamp >= $1 AND t.timestamp < $2
		 WHERE n.cultivo_id = $3
		 GROUP BY n.node_id, n.estado_id, e.estado_dela_conexion
		 ORDER BY n.node_id`,
		from,
		to,
		cultivoID,
	)
	if err != nil {
		return nil, err
	}
	defer nodoRows.Close()

	nodos := make([]models.NodoReport, 0)
	for nodoRows.Next() {
		var nodo models.NodoReport
		var ultima sql.NullTime
		if err := nodoRows.Scan(&nodo.NodeID, &nodo.EstadoID, &nodo.EstadoConexion, &ultima); err != nil {
			return nil, err
		}
		if ultima.Valid {
			value := ultima.Time
			nodo.UltimaLectura = &value
			report.NodosConDatos++
		}
		nodos = append(nodos, nodo)
	}
	if err := nodoRows.Err(); err != nil {
		return nil, err
	}

	report.Nodos = nodos
	report.NodosCount = len(nodos)

	var (
		avgSuelo sql.NullFloat64
		avgTemp  sql.NullFloat64
		avgAmb   sql.NullFloat64
		rows     sql.NullInt64
		ultima   sql.NullTime
	)
	if err := DB.QueryRow(
		queryCtx,
		`SELECT AVG(t.humedad_suelo), AVG(t.temperatura_amb), AVG(t.humedad_amb),
			COUNT(*) AS total_rows, MAX(t.timestamp)
		 FROM telemetria t
		 JOIN nodos n ON n.node_id = t.node_id
		 WHERE n.cultivo_id = $1 AND t.timestamp >= $2 AND t.timestamp < $3`,
		cultivoID,
		from,
		to,
	).Scan(&avgSuelo, &avgTemp, &avgAmb, &rows, &ultima); err != nil {
		return nil, err
	}

	if avgSuelo.Valid {
		value := avgSuelo.Float64
		report.AvgHumedadSuelo = &value
	}
	if avgTemp.Valid {
		value := avgTemp.Float64
		report.AvgTemperaturaAmb = &value
	}
	if avgAmb.Valid {
		value := avgAmb.Float64
		report.AvgHumedadAmb = &value
	}
	if rows.Valid {
		report.TelemetriaRows = rows.Int64
	}
	if ultima.Valid {
		value := ultima.Time
		report.UltimaLectura = &value
	}

	return &report, nil
}
