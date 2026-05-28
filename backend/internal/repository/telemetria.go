package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ardeuo_backend/internal/models"
)

const telemetriaQueryTimeout = 10 * time.Second

func ListTelemetria(ctx context.Context, limit *int, offset *int, start *time.Time, end *time.Time) ([]models.Telemetria, error) {
	queryCtx, cancel := context.WithTimeout(ctx, telemetriaQueryTimeout)
	defer cancel()

	query := `SELECT node_id, timestamp, humedad_suelo, temperatura_amb, humedad_amb
		FROM telemetria`
	args := make([]any, 0, 4)
	clauses := make([]string, 0, 1)
	argIndex := 1
	if start != nil && end != nil {
		clauses = append(clauses, fmt.Sprintf("timestamp >= $%d AND timestamp < $%d", argIndex, argIndex+1))
		args = append(args, *start, *end)
		argIndex += 2
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}
	query += " ORDER BY timestamp DESC, node_id"
	if limit != nil && *limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, *limit)
		argIndex++
	}
	if offset != nil && *offset >= 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, *offset)
	}

	rows, err := DB.Query(queryCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Telemetria, 0)
	for rows.Next() {
		var item models.Telemetria
		if err := rows.Scan(
			&item.NodeID,
			&item.Timestamp,
			&item.HumedadSuelo,
			&item.TemperaturaAmb,
			&item.HumedadAmb,
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
