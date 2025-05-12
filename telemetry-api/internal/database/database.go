package database

import (
	"context"
	"database/sql"
	"fmt"
	"telemetry-api/internal/models"
	"time"

	"telemetry-api/internal/observability"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(host string, port int, user, password, dbname string) (*Database, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) GetTelemetry(query *models.TelemetryQuery) ([]models.TelemetryRecord, int, error) {
	fmt.Println(query)

	baseQuery := `
		WITH results AS (
			SELECT timestamp, subsystem_id, temperature, battery, altitude, signal, has_anomaly,
				   COUNT(*) OVER() as total_count
			FROM telemetry
			WHERE timestamp BETWEEN $1 AND $2 AND subsystem_id = $3`

	offset := (query.Page - 1) * query.PageSize
	baseQuery += " ORDER BY timestamp DESC LIMIT $4 OFFSET $5"

	baseQuery += ")"
	baseQuery += `
		SELECT timestamp, subsystem_id, temperature, battery, altitude, signal, has_anomaly, total_count
		FROM results`

	rows, err := d.db.Query(baseQuery, query.StartTime, query.EndTime, query.SubsystemID, query.PageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying telemetry: %v", err)
	}
	defer rows.Close()

	var records []models.TelemetryRecord
	var totalCount int
	for rows.Next() {
		var record models.TelemetryRecord
		err := rows.Scan(
			&record.Timestamp,
			&record.SubsystemID,
			&record.Temperature,
			&record.Battery,
			&record.Altitude,
			&record.Signal,
			&record.HasAnomaly,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning telemetry record: %v", err)
		}
		records = append(records, record)
	}

	return records, totalCount, nil
}

func (d *Database) GetCurrentTelemetry() (*models.TelemetryRecord, error) {
	ctx := context.Background()
	start := time.Now()

	query := `
		SELECT timestamp, subsystem_id, temperature, battery, altitude, signal, has_anomaly
		FROM telemetry
		ORDER BY timestamp DESC
		LIMIT 1`

	var record models.TelemetryRecord
	err := d.db.QueryRowContext(ctx, query).Scan(
		&record.Timestamp,
		&record.SubsystemID,
		&record.Temperature,
		&record.Battery,
		&record.Altitude,
		&record.Signal,
		&record.HasAnomaly,
	)

	// Record metrics
	observability.RecordDBQuery(ctx, "get_current_telemetry", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("error getting current telemetry: %v", err)
	}

	return &record, nil
}

func (d *Database) GetAnomalies(query *models.TelemetryQuery) ([]models.AnomalyRecord, int, error) {
	ctx := context.Background()
	start := time.Now()

	baseQuery := `
		WITH results AS (
			SELECT timestamp, subsystem_id, anomaly_type, value, expected_range,
				   COUNT(*) OVER() as total_count
			FROM anomalies
			WHERE timestamp BETWEEN $1 AND $2`

	args := []interface{}{query.StartTime, query.EndTime}
	if query.SubsystemID != nil {
		baseQuery += " AND subsystem_id = $3"
		args = append(args, *query.SubsystemID)
	}

	baseQuery += " ORDER BY timestamp DESC LIMIT $4 OFFSET $5"
	offset := (query.Page - 1) * query.PageSize
	args = append(args, query.PageSize, offset)

	baseQuery += ")"
	baseQuery += `
		SELECT timestamp, subsystem_id, anomaly_type, value, expected_range, total_count
		FROM results`

	rows, err := d.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		observability.RecordDBQuery(ctx, "get_anomalies", time.Since(start), err)
		return nil, 0, fmt.Errorf("error querying anomalies: %v", err)
	}
	defer rows.Close()

	var records []models.AnomalyRecord
	var totalCount int
	for rows.Next() {
		var record models.AnomalyRecord
		err := rows.Scan(
			&record.Timestamp,
			&record.SubsystemID,
			&record.AnomalyType,
			&record.Value,
			&record.ExpectedRange,
			&totalCount,
		)
		if err != nil {
			observability.RecordDBQuery(ctx, "get_anomalies_scan", time.Since(start), err)
			return nil, 0, fmt.Errorf("error scanning anomaly record: %v", err)
		}
		records = append(records, record)
	}

	// Record successful query metrics
	observability.RecordDBQuery(ctx, "get_anomalies", time.Since(start), nil)

	return records, totalCount, nil
}

func (d *Database) GetAggregatedTelemetry(query *models.TelemetryAggregationQuery) ([]models.AggregatedMetric, error) {
	timeInterval := fmt.Sprintf("time_bucket('%s', timestamp)", query.GroupBy)

	sqlQuery := fmt.Sprintf(`
		SELECT 
			%s as timestamp,
			subsystem_id,
			MIN(temperature) as min,
			MAX(temperature) as max,
			AVG(temperature) as avg,
			COUNT(*) as count
		FROM telemetry
		WHERE timestamp BETWEEN $1 AND $2`, timeInterval)

	args := []interface{}{query.StartTime, query.EndTime}
	if query.SubsystemID != nil {
		sqlQuery += " AND subsystem_id = $3"
		args = append(args, *query.SubsystemID)
	}

	sqlQuery += fmt.Sprintf(" GROUP BY %s, subsystem_id ORDER BY timestamp DESC", timeInterval)

	rows, err := d.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying aggregated telemetry: %v", err)
	}
	defer rows.Close()

	var metrics []models.AggregatedMetric
	for rows.Next() {
		var metric models.AggregatedMetric
		err := rows.Scan(
			&metric.Timestamp,
			&metric.SubsystemID,
			&metric.Min,
			&metric.Max,
			&metric.Avg,
			&metric.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning aggregated metric: %v", err)
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
