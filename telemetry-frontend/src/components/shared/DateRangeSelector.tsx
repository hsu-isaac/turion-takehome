import React from "react";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { DateTimePicker } from "@mui/x-date-pickers/DateTimePicker";

interface DateRangeSelectorProps {
  startDate: Date | null;
  endDate: Date | null;
  onStartDateChange: (date: Date | null) => void;
  onEndDateChange: (date: Date | null) => void;
  onQuickSelect: (hours: number) => void;
  onRefresh: () => void;
  dateError: string | null;
  isLoading?: boolean;
}

const DateRangeSelector: React.FC<DateRangeSelectorProps> = ({
  startDate,
  endDate,
  onStartDateChange,
  onEndDateChange,
  onQuickSelect,
  onRefresh,
  dateError,
  isLoading = false,
}) => {
  return (
    <div className="mb-6 flex gap-4 items-center flex-wrap">
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <div className="flex flex-col">
          <DateTimePicker
            label="Start Date"
            value={startDate}
            onChange={onStartDateChange}
            maxDateTime={endDate || undefined}
          />
        </div>
        <div className="flex flex-col">
          <DateTimePicker
            label="End Date"
            value={endDate}
            onChange={onEndDateChange}
            minDateTime={startDate || undefined}
          />
        </div>
      </LocalizationProvider>
      {dateError && (
        <div className="text-red-500 text-sm mt-1">{dateError}</div>
      )}
      <div className="flex gap-2">
        <button
          onClick={() => onQuickSelect(6)}
          className="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
        >
          Last 6 Hours
        </button>
        <button
          onClick={() => onQuickSelect(24)}
          className="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
        >
          Last 24 Hours
        </button>
        <button
          onClick={onRefresh}
          disabled={isLoading}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? "Refreshing..." : "Refresh"}
        </button>
      </div>
    </div>
  );
};

export default DateRangeSelector;
