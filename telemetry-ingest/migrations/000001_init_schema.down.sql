-- Drop trigger and function
DROP TRIGGER IF EXISTS update_telemetry_updated_at ON telemetry;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indices
DROP INDEX IF EXISTS idx_telemetry_subsystem_timestamp;
DROP INDEX IF EXISTS idx_telemetry_has_anomaly_timestamp;
DROP INDEX IF EXISTS idx_anomalies_subsystem_timestamp;
DROP INDEX IF EXISTS idx_anomalies_type_timestamp;

-- Drop hypertables
DROP TABLE IF EXISTS telemetry CASCADE;
DROP TABLE IF EXISTS anomalies CASCADE;

-- Drop enum type
DROP TYPE IF EXISTS anomaly_type;

-- Note: We don't drop the timescaledb extension as it might be used by other tables 