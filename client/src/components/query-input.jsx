import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Table } from "@/components/ui/table";
import { unwrapAndDecodeData } from "@/utils.js";

// const keywords = ["READ", "WRITE", "DELETE", "family", "key", "qualifier", "value", "latest"];
// const operations = ["READ", "WRITE", "DELETE"];

  function RowKeyTables({ data }) {
    return (
      <>
        {Object.keys(data).map((rowKey) => {
          const row = data[rowKey];
          const qualifiers = Object.keys(row.cols).flatMap((family) =>
            Object.keys(row.cols[family])
          );

          return (
            <div key={rowKey} className="mt-4">
              <h3 className="text-md font-bold">Row Key: {rowKey}</h3>
              <Table className="table-auto border-collapse border border-gray-300 w-full text-sm">
                <thead>
                <tr>
                  <th className="border border-gray-300 px-2 py-1">Family</th>
                  {qualifiers.map((qualifier) => (
                    <th
                      key={qualifier}
                      className="border border-gray-300 px-2 py-1"
                    >
                      {qualifier}
                    </th>
                  ))}
                </tr>
                </thead>
                <tbody>
                {Object.keys(row.cols).map((family) => (
                  <tr key={family}>
                    <td className="border border-gray-300 px-2 py-1">
                      {family}
                    </td>
                    {qualifiers.map((qualifier) => (
                      <td
                        key={qualifier}
                        className="border border-gray-300 px-2 py-1"
                      >
                        {row.cols[family][qualifier]
                          ? row.cols[family][qualifier]
                            .map((item) => item.value)
                            .join(", ")
                          : "-"}
                      </td>
                    ))}
                  </tr>
                ))}
                </tbody>
              </Table>
            </div>
          );
        })}
      </>
    );
  }

export default function QueryInput() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState({});

  const handleChange = (e) => {
    const input = e.target.value;
    setQuery(input);
  };

  const handleSubmit = async () => {
    const payload = {
      query: query,
    };

    let response;
    try {
      response = await fetch("/query", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });
    } catch (error) {
      alert("Error: " + error.message);
      return;
    }

    let body;
    if (response.ok) {
      body = await response.json();
      // Normalize data to an array of rows\
      const transformedData = unwrapAndDecodeData(body);
      setResults(transformedData);
    } else {
      alert("error");
    }
  };

  return (
    <>
      <div className="flex flex-col gap-2 relative pt-4 w-[650px]">
        <label htmlFor="query" className="text-sm font-medium text-gray-700">
          LiteTable Query Input
        </label>
        <Input
          id="query"
          value={query}
          onChange={handleChange}
          className={`p-2 block w-full rounded-md border-gray-300 shadow-sm focus:ring-blue-500 sm:text-sm`}
          placeholder="Enter your query here..."
        />
        <Button
          onClick={handleSubmit}
          className="mt-2 bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded-md px-4 py-2"
        >
          Submit
        </Button>
      </div>
      <div className="mt-4">
        <h2 className="text-lg font-semibold">Results:</h2>
        {/*{results && <pre>{JSON.stringify(results, null, 2)}</pre>}*/}
        {results && (
          <div className="mt-4">
            <RowKeyTables data={results} />
          </div>
        )}
      </div>

    </>
  );
}
