"use client";

import { useState, Fragment } from "react";
import { format } from "date-fns";
import { ChevronDown, ChevronRight, X } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";

export default function ResultsTable({ data }) {
  const [expandedRows, setExpandedRows] = useState({});
  const [dialogValue, setDialogValue] = useState(null);

  const toggleRow = (key) => {
    setExpandedRows((prev) => ({
      ...prev,
      [key]: !prev[key],
    }));
  };

  const openDialogWithValue = (val) => {
    if (val === null || val === undefined) return;
    setDialogValue(val);
  };

  const renderValue = (val) => {
    if (val === null || val === undefined) return "-";
    if (typeof val === "object") return JSON.stringify(val);
    return String(val);
  };

  const allQualifiers = new Set();
  Object.values(data).forEach((row) => {
    Object.values(row.cols).forEach((families) => {
      Object.keys(families).forEach((qualifier) => {
        allQualifiers.add(qualifier);
      });
    });
  });
  const qualifierColumns = Array.from(allQualifiers).sort();

  const formatTimestamp = (timestamp) => {
    try {
      return format(new Date(timestamp), "MMM d, yyyy h:mm a");
    } catch (e) {
      return timestamp;
    }
  };

  const getLatestValue = (row, family, qualifier) => {
    if (!row.cols[family] || !row.cols[family][qualifier]) return null;

    const values = [...row.cols[family][qualifier]].sort(
      (a, b) =>
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
    );
    return values[0];
  };

  return (
    <div className="rounded-md border overflow-hidden">
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="bg-muted/50">
              <TableHead className="w-[50px]" />
              <TableHead className="w-[150px]">RowKey</TableHead>
              {qualifierColumns.map((qualifier) => (
                <TableHead key={qualifier} className="max-w-[200px]">
                  {qualifier}
                </TableHead>
              ))}
            </TableRow>
          </TableHeader>
          <TableBody>
            {Object.entries(data).map(([index, row]) => {
              const isExpanded = expandedRows[index] || false;
              const families = Object.keys(row.cols);

              return (
                <Fragment key={index}>
                  <TableRow
                    className="hover:bg-muted/50 cursor-pointer"
                    onClick={() => toggleRow(index)}
                  >
                    <TableCell className="p-2 text-center">
                      {isExpanded ? (
                        <ChevronDown className="h-4 w-4 inline" />
                      ) : (
                        <ChevronRight className="h-4 w-4 inline" />
                      )}
                    </TableCell>
                    <TableCell className="font-medium">{row.key}</TableCell>
                    {qualifierColumns.map((qualifier) => {
                      for (const family of families) {
                        const latestValue = getLatestValue(
                          row,
                          family,
                          qualifier,
                        );
                        if (latestValue) {
                          return (
                            <TableCell
                              key={qualifier}
                              className="max-w-[200px] truncate cursor-pointer"
                              onClick={() =>
                                openDialogWithValue(latestValue.value)
                              }
                            >
                              {renderValue(latestValue.value)}
                            </TableCell>
                          );
                        }
                      }
                      return <TableCell key={qualifier}>-</TableCell>;
                    })}
                  </TableRow>

                  {isExpanded && (
                    <TableRow key={`expanded-${index}`} className="bg-muted/20">
                      <TableCell
                        colSpan={2 + qualifierColumns.length}
                        className="p-0"
                      >
                        <div className="p-4">
                          <h4 className="text-sm font-medium mb-2">
                            Column Families
                          </h4>
                          {families.map((family) => (
                            <div key={family} className="mb-4">
                              <div className="flex items-center mb-2">
                                <Badge variant="outline" className="mr-2">
                                  {family}
                                </Badge>
                              </div>
                              <div className="overflow-x-auto">
                                <Table>
                                  <TableHeader>
                                    <TableRow className="bg-muted/30">
                                      {qualifierColumns.map((qualifier) =>
                                        row.cols[family] &&
                                        row.cols[family][qualifier] ? (
                                          <TableHead key={qualifier}>
                                            {qualifier}
                                          </TableHead>
                                        ) : null,
                                      )}
                                      <TableHead className="sticky left-0 bg-muted/30 z-10">
                                        Timestamp
                                      </TableHead>
                                    </TableRow>
                                  </TableHeader>
                                  <TableBody>
                                    {(() => {
                                      const allTimestamps = new Set();
                                      qualifierColumns.forEach((qualifier) => {
                                        if (
                                          row.cols[family] &&
                                          row.cols[family][qualifier]
                                        ) {
                                          row.cols[family][qualifier].forEach(
                                            (item) => {
                                              allTimestamps.add(item.timestamp);
                                            },
                                          );
                                        }
                                      });
                                      const sortedTimestamps = Array.from(
                                        allTimestamps,
                                      ).sort(
                                        (a, b) =>
                                          new Date(b).getTime() -
                                          new Date(a).getTime(),
                                      );
                                      return sortedTimestamps.map(
                                        (timestamp) => (
                                          <TableRow key={timestamp}>
                                            {qualifierColumns.map(
                                              (qualifier) => {
                                                if (
                                                  !row.cols[family] ||
                                                  !row.cols[family][qualifier]
                                                )
                                                  return null;
                                                const valueAtTimestamp =
                                                  row.cols[family][
                                                    qualifier
                                                  ].find(
                                                    (item) =>
                                                      item.timestamp ===
                                                      timestamp,
                                                  );
                                                return (
                                                  <TableCell
                                                    key={`${qualifier}-${timestamp}`}
                                                    className="max-w-[200px] truncate cursor-pointer"
                                                    onClick={() =>
                                                      openDialogWithValue(
                                                        valueAtTimestamp?.value,
                                                      )
                                                    }
                                                  >
                                                    {valueAtTimestamp
                                                      ? renderValue(
                                                          valueAtTimestamp.value,
                                                        )
                                                      : "-"}
                                                  </TableCell>
                                                );
                                              },
                                            )}
                                            <TableCell className="sticky left-0 z-10 font-medium">
                                              {formatTimestamp(timestamp)}
                                            </TableCell>
                                          </TableRow>
                                        ),
                                      );
                                    })()}
                                  </TableBody>
                                </Table>
                              </div>
                            </div>
                          ))}
                        </div>
                      </TableCell>
                    </TableRow>
                  )}
                </Fragment>
              );
            })}
          </TableBody>
        </Table>
      </div>

      {dialogValue !== null && (
        <div className="fixed inset-0 z-50 bg-black/40 flex items-center justify-center">
          <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-lg p-6 sm:max-w-[90vw] md:max-w-xl md w-full relative">
            <h2 className="text-lg font-semibold mb-4">Cell Value</h2>
            <pre className="bg-muted p-4 rounded text-sm overflow-auto max-h-[60vh] whitespace-pre-wrap break-words">
              {typeof dialogValue === "object"
                ? JSON.stringify(dialogValue, null, 2)
                : String(dialogValue)}
            </pre>
            <button
              className="absolute top-2 right-2 text-sm text-muted-foreground hover:text-foreground"
              onClick={() => setDialogValue(null)}
            >
              <X className="cursor-pointer" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
