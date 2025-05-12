"use client"

import { useState, Fragment } from "react"
import { format } from "date-fns"
import { ChevronDown, ChevronRight } from "lucide-react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"

export default function ResultsTable({ data }) {
	const [expandedRows, setExpandedRows] = useState({})

	const toggleRow = (key) => {
		setExpandedRows((prev) => ({
			...prev,
			[key]: !prev[key],
		}))
	}

	// Extract all unique qualifiers across all rows and sort alphabetically
	const allQualifiers = new Set()
	Object.values(data).forEach((row) => {
		Object.values(row.cols).forEach((families) => {
			Object.keys(families).forEach((qualifier) => {
				allQualifiers.add(qualifier)
			})
		})
	})

	// Convert to sorted array
	const qualifierColumns = Array.from(allQualifiers).sort()

	// Format timestamp for display
	const formatTimestamp = (timestamp) => {
		try {
			return format(new Date(timestamp), "MMM d, yyyy h:mm a")
		} catch (e) {
			return timestamp
		}
	}

	// Get the most recent value for a qualifier
	const getLatestValue = (row, family, qualifier) => {
		if (!row.cols[family] || !row.cols[family][qualifier]) return null

		// Sort by timestamp (newest first) and take the first one
		const values = [...row.cols[family][qualifier]].sort(
			(a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
		)

		return values[0]
	}

	return (
		<div className="rounded-md border overflow-hidden">
			<div className="overflow-x-auto">
				<Table>
					<TableHeader>
						<TableRow className="bg-muted/50">
							<TableHead className="w-[50px]"></TableHead>
							<TableHead className="w-[150px]">RowKey</TableHead>
							{qualifierColumns.map((qualifier) => (
								<TableHead key={qualifier}>{qualifier}</TableHead>
							))}
						</TableRow>
					</TableHeader>
					<TableBody>
						{Object.entries(data).map(([index, row]) => {
							const isExpanded = expandedRows[index] || false
							const families = Object.keys(row.cols)

							return (
								<Fragment key={index}>
									<TableRow className="hover:bg-muted/50 cursor-pointer" onClick={() => toggleRow(index)}>
										<TableCell className="p-2 text-center">
											{isExpanded ? (
												<ChevronDown className="h-4 w-4 inline" />
											) : (
												<ChevronRight className="h-4 w-4 inline" />
											)}
										</TableCell>
										<TableCell className="font-medium">{row.key}</TableCell>

										{qualifierColumns.map((qualifier) => {
											// Find this qualifier in any family
											for (const family of families) {
												const latestValue = getLatestValue(row, family, qualifier)
												if (latestValue) {
													return (
														<TableCell key={qualifier}>
															{typeof latestValue.value === "object"
																? JSON.stringify(latestValue.value)
																: latestValue.value}
														</TableCell>
													)
												}
											}
											return <TableCell key={qualifier}>-</TableCell>
										})}
									</TableRow>

									{isExpanded && (
										<TableRow key={`expanded-${index}`} className="bg-muted/20">
											<TableCell colSpan={2 + qualifierColumns.length} className="p-0">
												<div className="p-4">
													<h4 className="text-sm font-medium mb-2">Column Families</h4>
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
																				row.cols[family] && row.cols[family][qualifier] ? (
																					<TableHead key={qualifier}>{qualifier}</TableHead>
																				) : null,
																			)}
																			<TableHead className="sticky left-0 bg-muted/30 z-10">Timestamp</TableHead>
																		</TableRow>
																	</TableHeader>
																	<TableBody>
																		{(() => {
																			// Get all timestamps from all qualifiers in this family
																			const allTimestamps = new Set()

																			// Only collect timestamps for qualifiers that exist in this family
																			qualifierColumns.forEach((qualifier) => {
																				if (row.cols[family] && row.cols[family][qualifier]) {
																					row.cols[family][qualifier].forEach((item) => {
																						allTimestamps.add(item.timestamp)
																					})
																				}
																			})

																			// Sort timestamps (newest first)
																			const sortedTimestamps = Array.from(allTimestamps).sort(
																				(a, b) => new Date(b).getTime() - new Date(a).getTime(),
																			)

																			// For each timestamp, create a row with values for each qualifier
																			return sortedTimestamps.map((timestamp) => (
																				<TableRow key={timestamp}>
																					{qualifierColumns.map((qualifier) => {
																						// Skip qualifiers that don't exist in this family
																						if (!row.cols[family] || !row.cols[family][qualifier]) {
																							return null
																						}

																						// Find value for this qualifier at this timestamp
																						const valueAtTimestamp = row.cols[family][qualifier].find(
																							(item) => item.timestamp === timestamp,
																						)

																						return (
																							<TableCell key={`${qualifier}-${timestamp}`}>
																								{valueAtTimestamp
																									? typeof valueAtTimestamp.value === "object"
																										? JSON.stringify(valueAtTimestamp.value)
																										: valueAtTimestamp.value
																									: "-"}
																							</TableCell>
																						)
																					})}
																					<TableCell className="sticky left-0  z-10 font-medium">
																						{formatTimestamp(timestamp)}
																					</TableCell>
																				</TableRow>
																			))
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
							)
						})}
					</TableBody>
				</Table>
			</div>
		</div>
	)
}
