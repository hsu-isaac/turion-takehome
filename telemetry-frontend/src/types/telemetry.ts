export interface TelemetryRecord {
  timestamp: string;
  subsystem_id: number;
  temperature: number;
  battery: number;
  altitude: number;
  signal: number;
  has_anomaly: boolean;
}

export interface AnomalyRecord {
  timestamp: string;
  subsystem_id: number;
  anomaly_type: string;
  value: number;
  expected_range: string;
}

export interface TelemetryResponse {
  data: TelemetryRecord[] | AnomalyRecord[];
  metadata: {
    total_count: number;
    page_count: number;
    has_more: boolean;
    time_range: {
      start: string;
      end: string;
    };
  };
}

export interface TelemetryMetric {
  name: string;
  value: number;
  unit: string;
  status: "normal" | "warning" | "critical";
  range: {
    min: number;
    max: number;
  };
}
