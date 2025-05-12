"use client"

import { useState } from "react"
import { Plus, Trash2, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
	DialogFooter,
	DialogClose,
} from "@/components/ui/dialog"
import { toast } from "sonner"

import ResultsTable from "./results-table"
import {unwrapAndDecodeData} from '@/utils.js';

export default function QueryBuilder() {
	const [operation, setOperation] = useState("READ")
	const [filterType, setFilterType] = useState("key")
	const [filterValue, setFilterValue] = useState("")
	const [columnFamily, setColumnFamily] = useState("wrestlers")
	const [columnFamilies, setColumnFamilies] = useState(["wrestlers", "matches", "galaxies"])
	const [newColumnFamily, setNewColumnFamily] = useState("")
	const [latest, setLatest] = useState("1")
	const [qualifiers, setQualifiers] = useState([{ qualifier: "", value: "" }])
	const [generatedQuery, setGeneratedQuery] = useState("")
	const [results, setResults] = useState({});

	const [isLoading, setIsLoading] = useState(false)
	const [isOpen, setIsOpen] = useState("builder")

	const handleQualifierChange = (index, field, value) => {
		const updatedQualifiers = [...qualifiers]
		updatedQualifiers[index][field] = value
		setQualifiers(updatedQualifiers)
	}

	const addQualifier = () => {
		setQualifiers([...qualifiers, { qualifier: "", value: "" }])
	}

	const removeQualifier = (index) => {
		const updatedQualifiers = qualifiers.filter((_, i) => i !== index)
		setQualifiers(updatedQualifiers.length ? updatedQualifiers : [{ qualifier: "", value: "" }])
	}


	const buildQuery = () => {
		setIsLoading(true);
		setIsOpen("results");

		const payload = {
			type: operation,
			readType: filterType,
			key: filterValue || "",
			family: columnFamily || "",
			qualifiers: qualifiers
				.filter((q) => q.qualifier) // Only include qualifiers with a name
				.map((q) => {
					// For READ and DELETE operations, send just the qualifier name
					if (operation === "READ" || operation === "DELETE") {
						return { name: q.qualifier }; // Use 'name' property to match server expectations
					}
					// For WRITE operations, include both qualifier and value
					return {
						name: q.qualifier,
						value: encodeURIComponent(q.value || "")
					};
				}),
		};

		// Simulate loading
		setTimeout(() => {
			setGeneratedQuery(JSON.stringify(payload, null, 2)); // Display the payload as a JSON string
			setIsLoading(false);
		}, 1000);

		return payload;
	};

	const handleSubmit = async (payload) => {
		let response;
		try {
			response = await fetch("/query", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(payload),
			});
		} catch (e) {
			console.error("Error:", e);
			return;
		}

		let body;
		if (response.ok) {
			body = await response.json();
			// Normalize data to an array of rows\
			const transformedData = unwrapAndDecodeData(body);
			setResults(transformedData);
			if (operation === "WRITE") {
				toast("WRITE successful");
			} else if (operation === "READ") {
				toast("READ successful");
			} else if (operation === "DELETE") {
				toast("DELETE successful");
			}

		} else {
			console.error("Error:", response.json());
		}

		return body;
	};

	const addNewColumnFamily = async () => {
		const newest = newColumnFamily.trim()
		if (newColumnFamily && !columnFamilies.includes(newColumnFamily)) {
			setColumnFamilies([...columnFamilies, newColumnFamily])
			setColumnFamily(newest)
			setNewColumnFamily("")
		}

		toast("Column family created")
		const famQuery = {
			type: "CREATE",
			families: [newest],
		}

		await handleSubmit(famQuery)
	}


	const generate = async() => {
		const result = buildQuery();
		await handleSubmit(result);
	}


	return (
		<div className="space-y-4 mb-4">
			<Accordion type="single" value={isOpen} onValueChange={setIsOpen} collapsible className="w-full">
				<AccordionItem value="builder" className="border rounded-lg">
					<AccordionTrigger className="px-4 py-3 hover:no-underline">
						<span className="text-base font-medium">Query Builder</span>
					</AccordionTrigger>
					<AccordionContent className="px-6 pb-6 pt-2">
						<div className="grid grid-cols-2 gap-x-6 gap-y-4">
							<div>
								<Label htmlFor="operation" className="block mb-1">
									Operation
								</Label>
								<Select value={operation} onValueChange={(value) => setOperation(value)}>
									<SelectTrigger id="operation" className="w-full">
										<SelectValue placeholder="Select operation" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="READ">READ</SelectItem>
										<SelectItem value="WRITE">WRITE</SelectItem>
										<SelectItem value="DELETE">DELETE</SelectItem>
									</SelectContent>
								</Select>
							</div>

							<div>
								<div className="flex justify-between items-center mb-1">
									<Label htmlFor="columnFamily">Column Family</Label>
									<Dialog>
										<DialogTrigger asChild>
											<Button variant="ghost" size="sm" className="h-6 px-2">
												<Plus className="h-3.5 w-3.5 mr-1" />
												Add New
											</Button>
										</DialogTrigger>
										<DialogContent className="sm:max-w-md">
											<DialogHeader>
												<DialogTitle>Add New Column Family</DialogTitle>
											</DialogHeader>
											<div className="py-4">
												<Label htmlFor="newColumnFamily" className="mb-2 block">
													Name
												</Label>
												<Input
													id="newColumnFamily"
													value={newColumnFamily}
													onChange={(e) => setNewColumnFamily(e.target.value)}
													placeholder="Enter column family name"
												/>
											</div>
											<DialogFooter>
												<DialogClose asChild>
													<Button variant="outline">Cancel</Button>
												</DialogClose>
												<DialogClose asChild>
													<Button onClick={addNewColumnFamily} disabled={!newColumnFamily}>
														Add
													</Button>
												</DialogClose>
											</DialogFooter>
										</DialogContent>
									</Dialog>
								</div>
								<Select value={columnFamily} onValueChange={setColumnFamily}>
									<SelectTrigger id="columnFamily">
										<SelectValue placeholder="Select column family" />
									</SelectTrigger>
									<SelectContent>
										{columnFamilies.map((family) => (
											<SelectItem key={family} value={family}>
												{family}
											</SelectItem>
										))}
									</SelectContent>
								</Select>
							</div>

							<div>
								<Label htmlFor="filterType" className="block mb-1">
									Filter Type
								</Label>
								<Select value={filterType} onValueChange={(value) => setFilterType(value)}>
									<SelectTrigger id="filterType">
										<SelectValue placeholder="Select filter type" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="key">rowKey</SelectItem>
										<SelectItem value="prefix">prefix</SelectItem>
										<SelectItem value="regex">regex</SelectItem>
									</SelectContent>
								</Select>
							</div>

							<div>
								<Label htmlFor="filterValue" className="block mb-1">
									Filter Value
								</Label>
								<Input
									id="filterValue"
									placeholder="e.g. champ:1"
									value={filterValue}
									onChange={(e) => setFilterValue(e.target.value)}
								/>
							</div>

							{operation === "READ" && (
								<div>
									<Label htmlFor="latest" className="block mb-1">
										Latest
									</Label>
									<Input
										id="latest"
										type="number"
										placeholder="e.g. 3"
										value={latest}
										onChange={(e) => setLatest(e.target.value)}
										className="w-full"
									/>
								</div>
							)}
						</div>

						<div className="mt-4">
							<div className="flex justify-between items-center mb-2">
								<Label className="font-medium">Qualifiers</Label>
								<Button type="button" variant="outline" size="sm" onClick={addQualifier} className="h-8">
									<Plus className="h-4 w-4 mr-1" />
									Add
								</Button>
							</div>

							<div className="space-y-2">
								{qualifiers.map((pair, index) => (
									<div key={index} className="flex gap-2 items-center">
										<div className="flex-grow">
											<Input
												placeholder="Qualifier"
												value={pair.qualifier}
												onChange={(e) => handleQualifierChange(index, "qualifier", e.target.value)}
											/>
										</div>
										{operation === "WRITE" && (
											<div className="flex-grow">
												<Input
													placeholder="Value"
													value={pair.value}
													onChange={(e) => handleQualifierChange(index, "value", e.target.value)}
												/>
											</div>
										)}
										<Button
											type="button"
											variant="destructive"
											size="icon"
											onClick={() => removeQualifier(index)}
											className="shrink-0"
											disabled={qualifiers.length === 1}
										>
											<Trash2 className="h-4 w-4" />
											<span className="sr-only">Remove</span>
										</Button>
									</div>
								))}
							</div>
						</div>

						<Button onClick={generate} className="w-full mt-6 bg-black text-white hover:bg-black/90">
							Generate Query
						</Button>
					</AccordionContent>
				</AccordionItem>

				<AccordionItem value="results" className="border rounded-lg">
					<AccordionTrigger className="px-4 py-3 hover:no-underline">
						<span className="text-base font-medium">Results</span>
					</AccordionTrigger>
					<AccordionContent className="px-6 pb-6">
						{isLoading ? (
							<div className="flex justify-center items-center py-8">
								<Loader2 className="h-8 w-8 animate-spin text-primary" />
								<span className="ml-2">Processing query...</span>
							</div>
						) : generatedQuery ? (
							<div className="space-y-4">
								<div>
									<h3 className="text-sm font-medium mb-1">Generated Query:</h3>
									<div className="p-3 bg-muted rounded-md overflow-x-auto">
										<code className="text-sm">{generatedQuery}</code>
									</div>
								</div>

								{results && (
									<ResultsTable data={Object.values(results)} />
								)}
							</div>
						) : (
							<div className="py-4 text-center text-muted-foreground">Generate a query to see results</div>
						)}
					</AccordionContent>
				</AccordionItem>
			</Accordion>
		</div>
	)
}
