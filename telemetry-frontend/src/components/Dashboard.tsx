import React from "react";
import MetricCard from "./MetricCard";
import HistoricalData from "./HistoricalData";
import AnomaliesTable from "./AnomaliesTable";
import { useTelemetry } from "../hooks/useTelemetry";

const Dashboard: React.FC = () => {
  const { metrics, error } = useTelemetry();

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="text-xl text-red-600 bg-red-50 p-4 rounded-lg shadow-sm">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 space-y-8">
      <h1 className="text-3xl font-bold text-gray-900">
        Spacecraft Telemetry Dashboard
      </h1>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold text-gray-900">Real-time Updates</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {metrics.map((metric) => (
            <MetricCard key={metric.name} metric={metric} />
          ))}
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold text-gray-900">Anomalies</h2>
        <AnomaliesTable />
      </section>

      <section className="space-y-4">
        <HistoricalData />
      </section>
    </div>
  );
};

export default Dashboard;
