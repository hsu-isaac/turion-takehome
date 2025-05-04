package models

import "time"

type TelemetryRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	SubsystemID uint16    `json:"subsystem_id"`
	Temperature float32   `json:"temperature"`
	Battery     float32   `json:"battery"`
	Altitude    float32   `json:"altitude"`
	Signal      float32   `json:"signal"`
	HasAnomaly  bool      `json:"has_anomaly"`
}

type AnomalyRecord struct {
	Timestamp     time.Time `json:"timestamp"`
	SubsystemID   uint16    `json:"subsystem_id"`
	AnomalyType   string    `json:"anomaly_type"`
	Value         float32   `json:"value"`
	ExpectedRange string    `json:"expected_range"`
}

type TelemetryQuery struct {
	StartTime   time.Time `query:"start_time"`
	EndTime     time.Time `query:"end_time"`
	SubsystemID *uint16   `query:"subsystem_id"`
	Page        int       `query:"page" default:"1"`
	PageSize    int       `query:"page_size" default:"100"`
	Format      string    `query:"format" default:"raw"` // 'raw' or 'chart'
}

type TelemetryAggregationQuery struct {
	StartTime   time.Time `query:"start_time"`
	EndTime     time.Time `query:"end_time"`
	GroupBy     string    `query:"group_by"`    // '1m', '1h', '1d'
	Aggregation string    `query:"aggregation"` // 'min', 'max', 'avg', 'count'
	SubsystemID *uint16   `query:"subsystem_id"`
}

type TelemetryResponse struct {
	Data     interface{} `json:"data"`
	Metadata struct {
		TotalCount int       `json:"total_count"`
		PageCount  int       `json:"page_count"`
		HasMore    bool      `json:"has_more"`
		TimeRange  TimeRange `json:"time_range"`
	} `json:"metadata"`
}

type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type AggregatedMetric struct {
	Timestamp   time.Time `json:"timestamp"`
	SubsystemID uint16    `json:"subsystem_id"`
	Min         float32   `json:"min"`
	Max         float32   `json:"max"`
	Avg         float32   `json:"avg"`
	Count       int       `json:"count"`
}
