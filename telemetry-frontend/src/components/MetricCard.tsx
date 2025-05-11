import React from "react";
import { TelemetryMetric } from "../types/telemetry";

interface MetricCardProps {
  metric: TelemetryMetric;
}

const MetricCard: React.FC<MetricCardProps> = ({ metric }) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "critical":
        return "bg-red-50 border-red-200 text-red-700";
      case "warning":
        return "bg-yellow-50 border-yellow-200 text-yellow-700";
      default:
        return "bg-green-50 border-green-200 text-green-700";
    }
  };

  return (
    <div
      className={`rounded-lg border-2 p-4 shadow-sm transition-colors ${getStatusColor(
        metric.status
      )}`}
    >
      <h3 className="text-lg font-semibold mb-2">{metric.name}</h3>
      <div className="text-2xl font-bold mb-1">
        {metric.value.toFixed(2)} {metric.unit}
      </div>
      <div className="text-sm text-gray-600">
        Range: {metric.range.min} - {metric.range.max} {metric.unit}
      </div>
    </div>
  );
};

export default MetricCard;
