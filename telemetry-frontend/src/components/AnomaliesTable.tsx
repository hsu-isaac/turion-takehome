import React, { useEffect, useState, useMemo } from "react";
import { TelemetryService } from "../services/telemetryService";
import { AnomalyRecord, TelemetryResponse } from "../types/telemetry";
import DateRangeSelector from "./shared/DateRangeSelector";
import DataTable from "./shared/DataTable";

const ITEMS_PER_PAGE = 10;
const telemetryService = new TelemetryService();

const AnomaliesTable: React.FC = () => {
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
  const [anomalies, setAnomalies] = useState<AnomalyRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [pageCount, setPageCount] = useState(0);
  const [currentPageSort, setCurrentPageSort] = useState<{
    key: keyof AnomalyRecord;
    direction: "asc" | "desc";
  } | null>(null);

  const sortedData = useMemo(() => {
    if (!anomalies || anomalies.length === 0) {
      return [];
    }

    if (!currentPageSort) {
      return anomalies;
    }

    return [...anomalies].sort((a, b) => {
      if (currentPageSort.key === "timestamp") {
        const comparison =
          new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
        return currentPageSort.direction === "asc" ? comparison : -comparison;
      }

      if (currentPageSort.key === "value") {
        return currentPageSort.direction === "asc"
          ? a.value - b.value
          : b.value - a.value;
      }

      if (currentPageSort.key === "anomaly_type") {
        return currentPageSort.direction === "asc"
          ? a.anomaly_type.localeCompare(b.anomaly_type)
          : b.anomaly_type.localeCompare(a.anomaly_type);
      }

      return 0;
    });
  }, [anomalies, currentPageSort]);

  const fetchAnomalies = async () => {
    if (!startDate || !endDate) return;

    setLoading(true);
    setError(null);
    try {
      const response = await telemetryService.getAnomalies(
        startDate.toISOString(),
        endDate.toISOString(),
        currentPage,
        ITEMS_PER_PAGE
      );

      const anomalyData = response.data as AnomalyRecord[];
      setAnomalies(anomalyData);
      setTotalItems(response.metadata.total_count);
      setPageCount(response.metadata.page_count);
      setHasMore(response.metadata.has_more);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
      setAnomalies([]);
      setTotalItems(0);
      setPageCount(0);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    setCurrentPage(1);
    setCurrentPageSort(null);
    fetchAnomalies();
  }, [startDate, endDate]);

  useEffect(() => {
    fetchAnomalies();
  }, [currentPage]);

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

  const handleSort = (key: keyof AnomalyRecord) => {
    setCurrentPageSort((current) => ({
      key,
      direction:
        current?.key === key && current.direction === "asc" ? "desc" : "asc",
    }));
  };

  const columns = [
    {
      key: "timestamp" as keyof AnomalyRecord,
      header: "Timestamp",
      sortable: true,
      render: (value: string | number | boolean, item: AnomalyRecord) =>
        new Date(value as string).toLocaleString(),
    },
    {
      key: "anomaly_type" as keyof AnomalyRecord,
      header: "Metric",
      sortable: true,
      render: (value: string | number | boolean) => value as string,
    },
    {
      key: "value" as keyof AnomalyRecord,
      header: "Value",
      sortable: true,
      render: (value: string | number | boolean) =>
        (value as number).toFixed(2),
    },
    {
      key: "expected_range" as keyof AnomalyRecord,
      header: "Expected Range",
      sortable: false,
      render: (value: string | number | boolean) => value as string,
    },
  ];

  return (
    <div className="space-y-4">
      <DateRangeSelector
        startDate={startDate}
        endDate={endDate}
        onStartDateChange={handleStartDateChange}
        onEndDateChange={handleEndDateChange}
        onQuickSelect={handleQuickSelect}
        onRefresh={fetchAnomalies}
        dateError={dateError}
        isLoading={loading}
      />

      {error && (
        <div className="text-red-500 bg-red-50 p-4 rounded-md">{error}</div>
      )}

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
        isLoading={loading}
        error={error}
        emptyMessage="No anomalies detected for the selected time range"
      />
    </div>
  );
};

export default AnomaliesTable;
