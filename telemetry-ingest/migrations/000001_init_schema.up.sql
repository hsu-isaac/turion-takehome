-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Create enum for anomaly types
CREATE TYPE anomaly_type AS ENUM (
    'high_temperature',
    'low_temperature',
    'low_battery',
    'low_altitude',
    'weak_signal'
);

-- Create telemetry table
CREATE TABLE telemetry (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    subsystem_id SMALLINT NOT NULL,
    temperature FLOAT4 NOT NULL CHECK (temperature >= -273.15),
    battery FLOAT4 NOT NULL CHECK (battery >= 0 AND battery <= 100),
    altitude FLOAT4 NOT NULL CHECK (altitude >= 0),
    signal FLOAT4 NOT NULL,
    has_anomaly BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Create anomalies table
CREATE TABLE anomalies (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    subsystem_id SMALLINT NOT NULL,
    anomaly_type anomaly_type NOT NULL,
    value FLOAT4 NOT NULL,
    expected_range TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Convert tables to hypertables
SELECT create_hypertable('telemetry', 'timestamp',
    chunk_time_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

SELECT create_hypertable('anomalies', 'timestamp',
    chunk_time_interval => INTERVAL '1 hour',
    if_not_exists => TRUE
);

-- Create indices
CREATE INDEX IF NOT EXISTS idx_telemetry_subsystem_timestamp 
    ON telemetry (subsystem_id, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_telemetry_has_anomaly_timestamp 
    ON telemetry (has_anomaly, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_anomalies_subsystem_timestamp 
    ON anomalies (subsystem_id, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_anomalies_type_timestamp 
    ON anomalies (anomaly_type, timestamp DESC);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updated_at
CREATE TRIGGER update_telemetry_updated_at
    BEFORE UPDATE ON telemetry
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
