package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	errorCounter    metric.Int64Counter

	dbQueryCounter  metric.Int64Counter
	dbQueryDuration metric.Float64Histogram
	dbErrorCounter  metric.Int64Counter

	telemetryCounter metric.Int64Counter
	anomalyCounter   metric.Int64Counter
)

func InitializeMetrics() error {
	meter := otel.GetMeterProvider().Meter("telemetry-api")

	var err error

	requestCounter, err = meter.Int64Counter(
		"http.requests.total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	requestDuration, err = meter.Float64Histogram(
		"http.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	errorCounter, err = meter.Int64Counter(
		"http.errors.total",
		metric.WithDescription("Total number of HTTP errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	dbQueryCounter, err = meter.Int64Counter(
		"db.queries.total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	dbQueryDuration, err = meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	dbErrorCounter, err = meter.Int64Counter(
		"db.errors.total",
		metric.WithDescription("Total number of database errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	telemetryCounter, err = meter.Int64Counter(
		"telemetry.records.total",
		metric.WithDescription("Total number of telemetry records processed"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	anomalyCounter, err = meter.Int64Counter(
		"telemetry.anomalies.total",
		metric.WithDescription("Total number of anomalies detected"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	return nil
}

func RecordRequest(ctx context.Context, method, path string, status int, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.Int("status", status),
	}

	requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	if status >= 400 {
		errorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

func RecordDBQuery(ctx context.Context, queryType string, duration time.Duration, err error) {
	attrs := []attribute.KeyValue{
		attribute.String("query_type", queryType),
	}

	dbQueryCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	dbQueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	if err != nil {
		dbErrorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

func RecordTelemetry(ctx context.Context, hasAnomaly bool) {
	telemetryCounter.Add(ctx, 1)
	if hasAnomaly {
		anomalyCounter.Add(ctx, 1)
	}
}
