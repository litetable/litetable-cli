"use client";

import { useState } from "react";
import { Loader2, Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { useQueryApi } from "@/hooks/use-query-api.js";
import ResultsTable from "./results-table";
import FamilySelector from "@/components/families.jsx";

export default function QueryBuilder() {
  const [operation, setOperation] = useState("READ");
  const [filterType, setFilterType] = useState("key");
  const [filterValue, setFilterValue] = useState("");
  const [columnFamily, setColumnFamily] = useState("wrestlers");
  const [latest, setLatest] = useState("1");
  const [qualifiers, setQualifiers] = useState([{ qualifier: "", value: "" }]);
  const [generatedQuery, setGeneratedQuery] = useState("");

  const [isOpen, setIsOpen] = useState("builder");
  const { handleSubmit, results, isLoading } = useQueryApi();

  const handleQualifierChange = (index, field, value) => {
    const updatedQualifiers = [...qualifiers];
    updatedQualifiers[index][field] = value;
    setQualifiers(updatedQualifiers);
  };

  const addQualifier = () => {
    setQualifiers([...qualifiers, { qualifier: "", value: "" }]);
  };

  const removeQualifier = (index) => {
    const updatedQualifiers = qualifiers.filter((_, i) => i !== index);
    setQualifiers(
      updatedQualifiers.length
        ? updatedQualifiers
        : [{ qualifier: "", value: "" }],
    );
  };

  const buildQuery = () => {
    setIsOpen("results");
    const parsedLatest = parseInt(latest, 10);
    return {
      type: operation,
      readType: filterType,
      key: filterValue || "",
      family: columnFamily || "",
      latest:
        operation === "READ" ? (isNaN(parsedLatest) ? 1 : parsedLatest) : 0,
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
            value: encodeURIComponent(q.value || ""),
          };
        }),
    };
  };

  const generate = async () => {
    const payload = buildQuery();
    setGeneratedQuery(JSON.stringify(payload, null, 2));
    await handleSubmit(payload);
  };

  return (
    <div className="mb-4">
      <h1 className={"text-3xl font-bold text-center mb-6"}>Query Builder</h1>
      <Accordion
        type="single"
        value={isOpen}
        onValueChange={setIsOpen}
        collapsible
        className="w-full space-y-4"
      >
        <AccordionItem value="builder" className="border rounded-lg">
          <AccordionTrigger className="px-4 py-3 hover:no-underline">
            <span className="text-base font-medium">Query Builder</span>
          </AccordionTrigger>
          <AccordionContent className="px-6 pb-6 pt-2">
            <div className="grid grid-cols-2 gap-x-6 gap-y-4">
              <div className={"flex flex-col gap-2"}>
                <Label htmlFor="operation" className="block mb-1">
                  Operation
                </Label>
                <Select
                  value={operation}
                  onValueChange={(value) => setOperation(value)}
                >
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
              <FamilySelector
                columnFamily={columnFamily}
                setColumnFamily={setColumnFamily}
                handleSubmit={handleSubmit}
              />
              <div className={"flex justify-between"}>
                <div>
                  <Label htmlFor="filterType" className="block mb-1">
                    Filter Type
                  </Label>
                  <Select
                    value={filterType}
                    onValueChange={(value) => setFilterType(value)}
                  >
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
            </div>

            <div className="mt-4">
              <div className="flex justify-between items-center mb-2">
                <Label className="font-medium">Qualifiers</Label>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={addQualifier}
                  className="h-8"
                >
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
                        onChange={(e) =>
                          handleQualifierChange(
                            index,
                            "qualifier",
                            e.target.value,
                          )
                        }
                      />
                    </div>
                    {operation === "WRITE" && (
                      <div className="flex-grow">
                        <Input
                          placeholder="Value"
                          value={pair.value}
                          onChange={(e) =>
                            handleQualifierChange(
                              index,
                              "value",
                              e.target.value,
                            )
                          }
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

            <Button
              onClick={generate}
              className="w-full mt-6 bg-black text-white hover:bg-black/90"
            >
              Generate Query
            </Button>
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="results" className="border rounded-lg ">
          <AccordionTrigger className="px-4 py-3 hover:no-underline">
            <span className="text-base font-medium">Results</span>
          </AccordionTrigger>
          <AccordionContent className="px-6 pb-6 border-t">
            {isLoading ? (
              <div className="flex justify-center items-center py-4">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <span className="ml-2">Processing query...</span>
              </div>
            ) : generatedQuery ? (
              <div className="mt-4">
                {/*<div>*/}
                {/*  <h3 className="text-sm font-medium mb-1">Generated Query:</h3>*/}
                {/*  <div className="p-3 bg-muted rounded-md overflow-x-auto">*/}
                {/*    <code className="text-sm">{generatedQuery}</code>*/}
                {/*  </div>*/}
                {/*</div>*/}
                {results && <ResultsTable data={Object.values(results)} />}
              </div>
            ) : (
              <div className="py-2 text-center text-muted-foreground">
                Generate a query to see results
              </div>
            )}
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
}
