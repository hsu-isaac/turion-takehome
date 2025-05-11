import React from "react";

interface Column<T> {
  key: keyof T;
  header: string;
  render?: (value: T[keyof T], item: T) => React.ReactNode;
  sortable?: boolean;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  currentSort: { key: keyof T; direction: "asc" | "desc" } | null;
  onSort: (key: keyof T) => void;
  currentPage: number;
  pageCount: number;
  totalItems: number;
  itemsPerPage: number;
  hasMore: boolean;
  onPageChange: (page: number) => void;
  isLoading?: boolean;
  error?: string | null;
  emptyMessage?: string;
}

function DataTable<T>({
  columns,
  data,
  currentSort,
  onSort,
  currentPage,
  pageCount,
  totalItems,
  itemsPerPage,
  hasMore,
  onPageChange,
  isLoading = false,
  error = null,
  emptyMessage = "No data available",
}: DataTableProps<T>) {
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {columns.map((column) => (
                <th
                  key={String(column.key)}
                  className={`px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${
                    column.sortable ? "cursor-pointer hover:bg-gray-100" : ""
                  }`}
                  onClick={() => column.sortable && onSort(column.key)}
                >
                  {column.header}{" "}
                  {column.sortable && currentSort?.key === column.key && (
                    <span>{currentSort.direction === "asc" ? "↑" : "↓"}</span>
                  )}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {error ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-6 py-4 text-center text-red-500"
                >
                  {error}
                </td>
              </tr>
            ) : data.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length}
                  className="px-6 py-4 text-center text-gray-500"
                >
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              data.map((item, index) => (
                <tr key={index}>
                  {columns.map((column) => (
                    <td
                      key={String(column.key)}
                      className="px-6 py-4 whitespace-nowrap text-sm text-gray-900"
                    >
                      {column.render
                        ? column.render(item[column.key], item)
                        : String(item[column.key])}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <div className="flex justify-between items-center bg-white p-4 border-t">
        <div className="text-sm text-gray-700">
          {totalItems > 0 ? (
            <>
              Showing {(currentPage - 1) * itemsPerPage + 1} to{" "}
              {Math.min(currentPage * itemsPerPage, totalItems)} of {totalItems}{" "}
              entries
              {hasMore && " (more available)"}
            </>
          ) : (
            "No entries to display"
          )}
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => onPageChange(Math.max(1, currentPage - 1))}
            disabled={currentPage === 1 || isLoading}
            className="px-3 py-1 rounded border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Previous
          </button>
          <span className="px-3 py-1">
            Page {currentPage} of {pageCount || 1}
          </span>
          <button
            onClick={() => onPageChange(Math.min(pageCount, currentPage + 1))}
            disabled={currentPage === pageCount || isLoading || !hasMore}
            className="px-3 py-1 rounded border hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}

export default DataTable;
