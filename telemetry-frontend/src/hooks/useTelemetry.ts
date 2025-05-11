import { useState, useEffect, useCallback } from "react";
import { TelemetryService } from "../services/telemetryService";
import { TelemetryRecord, TelemetryMetric } from "../types/telemetry";

const telemetryService = new TelemetryService();

export const useTelemetry = () => {
  const [metrics, setMetrics] = useState<TelemetryMetric[]>([]);
  const [error, setError] = useState<string | null>(null);

  const processMetrics = useCallback(
    (data: TelemetryRecord): TelemetryMetric[] => {
      return [
        {
          name: "Temperature",
          value: data.temperature,
          unit: "°C",
          status:
            data.temperature > 35
              ? "critical"
              : data.temperature > 30
              ? "warning"
              : "normal",
          range: { min: 20, max: 30 },
        },
        {
          name: "Battery",
          value: data.battery,
          unit: "%",
          status:
            data.battery < 40
              ? "critical"
              : data.battery < 70
              ? "warning"
              : "normal",
          range: { min: 70, max: 100 },
        },
        {
          name: "Altitude",
          value: data.altitude,
          unit: "km",
          status:
            data.altitude < 400
              ? "critical"
              : data.altitude < 500
              ? "warning"
              : "normal",
          range: { min: 500, max: 550 },
        },
        {
          name: "Signal Strength",
          value: data.signal,
          unit: "dB",
          status:
            data.signal < -80
              ? "critical"
              : data.signal < -60
              ? "warning"
              : "normal",
          range: { min: -60, max: -40 },
        },
      ];
    },
    []
  );

  const handleTelemetryUpdate = useCallback(
    (data: TelemetryRecord) => {
      setMetrics(processMetrics(data));
    },
    [processMetrics]
  );

  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        const data = await telemetryService.getCurrentTelemetry();
        handleTelemetryUpdate(data);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to fetch telemetry data"
        );
      }
    };

    fetchInitialData();
    telemetryService.connectWebSocket(handleTelemetryUpdate);

    return () => {
      telemetryService.disconnectWebSocket();
    };
  }, [handleTelemetryUpdate]);

  return {
    metrics,
    error,
  };
};
