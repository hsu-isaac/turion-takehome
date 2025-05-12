import React, { useState, useEffect, useMemo, useCallback } from "react";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";
import { TelemetryService } from "../services/telemetryService";
import { TelemetryRecord } from "../types/telemetry";
import DateRangeSelector from "./shared/DateRangeSelector";
import DataTable from "./shared/DataTable";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend
);

const telemetryService = new TelemetryService();
const ITEMS_PER_PAGE = 20;

const HistoricalData: React.FC = () => {
  const getInitialDateRange = () => {
    const end = new Date();
    const start = new Date(end.getTime() - 6 * 60 * 60 * 1000);
    return { start, end };
  };

  const initialDateRange = getInitialDateRange();
  const [startDate, setStartDate] = useState<Date | null>(
    initialDateRange.start
  );
  const [endDate, setEndDate] = useState<Date | null>(initialDateRange.end);
  const [dateError, setDateError] = useState<string | null>(null);
  const [telemetryData, setTelemetryData] = useState<TelemetryRecord[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [pageCount, setPageCount] = useState(0);
  const [currentPageSort, setCurrentPageSort] = useState<{
    key: keyof TelemetryRecord;
    direction: "asc" | "desc";
  } | null>(null);

  const sortedData = useMemo(() => {
    if (!telemetryData || telemetryData.length === 0) {
      return [];
    }

    if (!currentPageSort) {
      return telemetryData;
    }

    return [...telemetryData].sort((a, b) => {
      if (currentPageSort.key === "timestamp") {
        const comparison =
          new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
        return currentPageSort.direction === "asc" ? comparison : -comparison;
      }

      const aValue = a[currentPageSort.key];
      const bValue = b[currentPageSort.key];

      if (typeof aValue === "number" && typeof bValue === "number") {
        return currentPageSort.direction === "asc"
          ? aValue - bValue
          : bValue - aValue;
      }

      return 0;
    });
  }, [telemetryData, currentPageSort]);

  const chartData = useMemo(() => {
    if (!telemetryData || telemetryData.length === 0) {
      return {
        labels: [],
        datasets: [
          {
            label: "Temperature (°C)",
            data: [],
            borderColor: "rgb(255, 99, 132)",
            tension: 0.1,
          },
          {
            label: "Battery (%)",
            data: [],
            borderColor: "rgb(54, 162, 235)",
            tension: 0.1,
          },
          {
            label: "Altitude (km)",
            data: [],
            borderColor: "rgb(75, 192, 192)",
            tension: 0.1,
          },
          {
            label: "Signal (dB)",
            data: [],
            borderColor: "rgb(153, 102, 255)",
            tension: 0.1,
          },
        ],
      };
    }

    return {
      labels: telemetryData.map((record) =>
        new Date(record.timestamp).toLocaleTimeString()
      ),
      datasets: [
        {
          label: "Temperature (°C)",
          data: telemetryData.map((record) => record.temperature),
          borderColor: "rgb(255, 99, 132)",
          tension: 0.1,
        },
        {
          label: "Battery (%)",
          data: telemetryData.map((record) => record.battery),
          borderColor: "rgb(54, 162, 235)",
          tension: 0.1,
        },
        {
          label: "Altitude (km)",
          data: telemetryData.map((record) => record.altitude),
          borderColor: "rgb(75, 192, 192)",
          tension: 0.1,
        },
        {
          label: "Signal (dB)",
          data: telemetryData.map((record) => record.signal),
          borderColor: "rgb(153, 102, 255)",
          tension: 0.1,
        },
      ],
    };
  }, [telemetryData]);

  const fetchData = useCallback(async () => {
    if (!startDate || !endDate) return;

    setIsLoading(true);
    setError(null);
    try {
      const response = await telemetryService.getTelemetryHistory(
        startDate.toISOString(),
        endDate.toISOString(),
        currentPage,
        ITEMS_PER_PAGE
      );

      if (!response || !response.data || !Array.isArray(response.data)) {
        throw new Error("Invalid response format from server");
      }

      setTelemetryData(response.data as TelemetryRecord[]);
      setTotalItems(response.metadata.total_count);
      setPageCount(response.metadata.page_count);
      setHasMore(response.metadata.has_more);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch data");
      setTelemetryData([]);
      setTotalItems(0);
      setPageCount(0);
      setHasMore(false);
    } finally {
      setIsLoading(false);
    }
  }, [startDate, endDate, currentPage]);

  useEffect(() => {
    setCurrentPage(1);
    setCurrentPageSort(null);
    fetchData();
  }, [startDate, endDate, fetchData]);

  useEffect(() => {
    fetchData();
  }, [currentPage, fetchData]);

  const handleSort = (key: keyof TelemetryRecord) => {
    setCurrentPageSort((current) => ({
      key,
      direction:
        current?.key === key && current.direction === "asc" ? "desc" : "asc",
    }));
  };

  const handleQuickSelect = (hours: number) => {
    const end = new Date();
    const start = new Date(end.getTime() - hours * 60 * 60 * 1000);
    setStartDate(start);
    setEndDate(end);
  };

  const handleStartDateChange = (newValue: Date | null) => {
    if (newValue && endDate && newValue > endDate) {
      setDateError("Start date cannot be after end date");
      return;
    }
    setDateError(null);
    setStartDate(newValue);
  };

  const handleEndDateChange = (newValue: Date | null) => {
    if (newValue && startDate && newValue < startDate) {
      setDateError("End date cannot be before start date");
      return;
    }
    setDateError(null);
    setEndDate(newValue);
  };

  const chartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: "top" as const,
      },
      title: {
        display: true,
        text: "Telemetry History",
      },
    },
    scales: {
      y: {
        beginAtZero: false,
      },
    },
  };

  const columns = [
    {
      key: "timestamp" as keyof TelemetryRecord,
      header: "Timestamp",
      sortable: true,
      render: (value: string | number | boolean, item: TelemetryRecord) =>
        new Date(value as string).toLocaleString(),
    },
    {
      key: "temperature" as keyof TelemetryRecord,
      header: "Temperature (°C)",
      sortable: true,
      render: (value: string | number | boolean, item: TelemetryRecord) =>
        (value as number).toFixed(2),
    },
    {
      key: "battery" as keyof TelemetryRecord,
      header: "Battery (%)",
      sortable: true,
      render: (value: string | number | boolean, item: TelemetryRecord) =>
        (value as number).toFixed(2),
    },
    {
      key: "altitude" as keyof TelemetryRecord,
      header: "Altitude (km)",
      sortable: true,
      render: (value: string | number | boolean, item: TelemetryRecord) =>
        (value as number).toFixed(2),
    },
    {
      key: "signal" as keyof TelemetryRecord,
      header: "Signal (dB)",
      sortable: true,
      render: (value: string | number | boolean, item: TelemetryRecord) =>
        (value as number).toFixed(2),
    },
    {
      key: "has_anomaly" as keyof TelemetryRecord,
      header: "Status",
      sortable: false,
      render: (value: string | number | boolean, item: TelemetryRecord) => (
        <span
          className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
            value ? "bg-red-100 text-red-800" : "bg-green-100 text-green-800"
          }`}
        >
          {value ? "Anomaly" : "Normal"}
        </span>
      ),
    },
  ];

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold mb-4">Historical Data</h2>

      <DateRangeSelector
        startDate={startDate}
        endDate={endDate}
        onStartDateChange={handleStartDateChange}
        onEndDateChange={handleEndDateChange}
        onQuickSelect={handleQuickSelect}
        onRefresh={fetchData}
        dateError={dateError}
        isLoading={isLoading}
      />

      {error && (
        <div className="text-red-500 bg-red-50 p-4 rounded-md">{error}</div>
      )}

      <div className="bg-white rounded-lg shadow p-4 mb-8">
        <h3 className="text-xl font-semibold mb-4">Telemetry History Graph</h3>
        <Line data={chartData} options={chartOptions} />
      </div>

      <DataTable
        columns={columns}
        data={sortedData}
        currentSort={currentPageSort}
        onSort={handleSort}
        currentPage={currentPage}
        pageCount={pageCount}
        totalItems={totalItems}
        itemsPerPage={ITEMS_PER_PAGE}
        hasMore={hasMore}
        onPageChange={setCurrentPage}
        isLoading={isLoading}
        error={error}
        emptyMessage="No data available for the selected time range"
      />
    </div>
  );
};

export default HistoricalData;
