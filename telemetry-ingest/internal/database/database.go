package database

import (
	"database/sql"
	"fmt"
	"telemetry-ingest/internal/models"

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

func (d *Database) InitializeTables() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS telemetry (
			timestamp TIMESTAMPTZ NOT NULL,
			subsystem_id SMALLINT NOT NULL,
			temperature FLOAT NOT NULL,
			battery FLOAT NOT NULL,
			altitude FLOAT NOT NULL,
			signal FLOAT NOT NULL,
			has_anomaly BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating telemetry table: %v", err)
	}

	_, err = d.db.Exec(`
		SELECT create_hypertable('telemetry', 'timestamp', 
			chunk_time_interval => INTERVAL '1 hour',
			if_not_exists => TRUE
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating hypertable: %v", err)
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS anomalies (
			timestamp TIMESTAMPTZ NOT NULL,
			subsystem_id SMALLINT NOT NULL,
			anomaly_type TEXT NOT NULL,
			value FLOAT NOT NULL,
			expected_range TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating anomalies table: %v", err)
	}

	_, err = d.db.Exec(`
		SELECT create_hypertable('anomalies', 'timestamp',
			chunk_time_interval => INTERVAL '1 hour',
			if_not_exists => TRUE
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating anomalies hypertable: %v", err)
	}

	return nil
}

func (d *Database) StoreTelemetry(record *models.TelemetryRecord) error {
	_, err := d.db.Exec(`
		INSERT INTO telemetry (
			timestamp, subsystem_id, temperature, battery, altitude, signal, has_anomaly
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		record.Timestamp, record.SubsystemID, record.Temperature,
		record.Battery, record.Altitude, record.Signal, record.HasAnomaly,
	)
	if err != nil {
		return fmt.Errorf("error storing telemetry: %v", err)
	}
	return nil
}

func (d *Database) StoreAnomalies(anomalies []models.Anomaly) error {
	if len(anomalies) == 0 {
		return nil
	}

	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO anomalies (
			timestamp, subsystem_id, anomaly_type, value, expected_range
		) VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, anomaly := range anomalies {
		_, err = stmt.Exec(
			anomaly.Timestamp,
			anomaly.SubsystemID,
			anomaly.AnomalyType,
			anomaly.Value,
			anomaly.ExpectedRange,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error storing anomaly: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
