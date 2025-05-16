import { useState } from "react";
import { Upload, FileText, AlertCircle, Loader2, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Progress } from "@/components/ui/progress";
import { useQueryApi } from "@/hooks/use-query-api";
import { toast } from "sonner";

export default function JsonUpload() {
  const [jsonContent, setJsonContent] = useState("");
  const [uploadStatus, setUploadStatus] = useState("idle"); // idle, validating, processing, success, error
  const [errorMessage, setErrorMessage] = useState("");
  const [progress, setProgress] = useState(0);
  const [processedCount, setProcessedCount] = useState(0);
  const [totalCount, setTotalCount] = useState(0);
  const { handleSubmit } = useQueryApi();

  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      setJsonContent(e.target.result);
      toast("File loaded successfully");
    };
    reader.onerror = () => {
      toast.error("Error reading file");
    };
    reader.readAsText(file);
  };

  const handleClear = () => {
    setJsonContent("");
    setUploadStatus("idle");
    setErrorMessage("");
    setProgress(0);
    setProcessedCount(0);
    setTotalCount(0);
  };

  const validateJson = (jsonData) => {
    if (!Array.isArray(jsonData)) {
      jsonData = [jsonData]; // Convert single object to array
    }

    for (let i = 0; i < jsonData.length; i++) {
      const item = jsonData[i];
      if (!item.rowkey)
        return { valid: false, message: `Item ${i + 1} is missing rowkey` };
      if (!item.family)
        return { valid: false, message: `Item ${i + 1} is missing family` };
      if (!item.qualifiers || typeof item.qualifiers !== "object") {
        return {
          valid: false,
          message: `Item ${i + 1} has invalid qualifiers`,
        };
      }
    }

    return { valid: true, data: jsonData };
  };

  const processJson = async () => {
    try {
      setUploadStatus("validating");
      setErrorMessage("");
      setProgress(0);
      setProcessedCount(0);

      // Parse JSON
      const jsonData = JSON.parse(jsonContent);

      // Validate JSON structure
      const validation = validateJson(jsonData);
      if (!validation.valid) {
        setUploadStatus("error");
        setErrorMessage(validation.message);
        return;
      }

      const dataArray = validation.data;
      setTotalCount(dataArray.length);
      setUploadStatus("processing");

      // Process each item
      for (let i = 0; i < dataArray.length; i++) {
        const item = dataArray[i];

        // Convert to payload format
        const payload = {
          query: "WRITE",
          type: "WRITE",
          key: item.rowkey,
          family: item.family,
          qualifiers: Object.entries(item.qualifiers).map(
            ([qualifier, value]) => ({
              name: qualifier,
              value: encodeURIComponent(value),
            }),
          ),
        };

        // Submit the query
        await handleSubmit(payload);

        // Update progress
        setProcessedCount(i + 1);
        setProgress(Math.floor(((i + 1) / dataArray.length) * 100));
      }

      setUploadStatus("success");
      toast.success(`Processed ${dataArray.length} items successfully`);
    } catch (error) {
      setUploadStatus("error");
      setErrorMessage(error.message || "Failed to process JSON");
      toast.error("Error processing JSON");
    }
  };

  return (
    <>
      <h1 className="text-3xl font-bold text-center mb-6">Upload</h1>
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Upload JSON Data</CardTitle>
          <CardDescription>
            Upload JSON file or paste JSON content to process multiple records
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="upload" className="w-full">
            <TabsList className="grid grid-cols-2 mb-4">
              <TabsTrigger value="upload">Upload File</TabsTrigger>
              <TabsTrigger value="paste">Paste JSON</TabsTrigger>
            </TabsList>

            <TabsContent value="upload" className="space-y-4">
              <div className="grid w-full max-w-sm items-center gap-1.5">
                <label
                  htmlFor="json-file"
                  className="cursor-pointer border-2 border-dashed rounded-md p-6 text-center hover:bg-muted/50"
                >
                  <Upload className="h-8 w-8 mx-auto mb-2 text-muted-foreground" />
                  <p className="text-sm font-medium">
                    Click to upload JSON file
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    or drag and drop
                  </p>
                  <input
                    id="json-file"
                    type="file"
                    accept=".json"
                    onChange={handleFileUpload}
                    className="hidden"
                  />
                </label>
              </div>
            </TabsContent>

            <TabsContent value="paste" className="space-y-4">
              <Textarea
                placeholder="Paste JSON content here..."
                className="min-h-[200px] max-h-[300px] overflow-y-auto resize-none"
                value={jsonContent}
                onChange={(e) => setJsonContent(e.target.value)}
              />
            </TabsContent>

            {jsonContent && (
              <div className="mt-4 space-y-4">
                <div className="flex justify-between items-center">
                  <p className="text-sm font-medium">JSON Content Preview</p>
                  <Button variant="ghost" size="sm" onClick={handleClear}>
                    Clear
                  </Button>
                </div>
                <div className="bg-muted rounded-md p-3 max-h-[200px] overflow-y-auto">
                  <pre className="text-xs">
                    {jsonContent.slice(0, 500)}
                    {jsonContent.length > 500 ? "..." : ""}
                  </pre>
                </div>
              </div>
            )}

            {uploadStatus === "error" && (
              <Alert variant="destructive" className="mt-4">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{errorMessage}</AlertDescription>
              </Alert>
            )}

            {uploadStatus === "processing" && (
              <div className="mt-4 space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Processing records...</span>
                  <span>
                    {processedCount} of {totalCount}
                  </span>
                </div>
                <Progress value={progress} className="h-2" />
              </div>
            )}

            {uploadStatus === "success" && (
              <Alert className="mt-4 bg-green-50 border-green-200">
                <Check className="h-4 w-4 text-green-600" />
                <AlertTitle className="text-green-600">Success</AlertTitle>
                <AlertDescription>
                  All {totalCount} records were processed successfully.
                </AlertDescription>
              </Alert>
            )}

            <Button
              onClick={processJson}
              disabled={!jsonContent || uploadStatus === "processing"}
              className="w-full mt-4"
            >
              {uploadStatus === "processing" ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Processing...
                </>
              ) : (
                <>
                  <FileText className="mr-2 h-4 w-4" />
                  Process JSON
                </>
              )}
            </Button>
          </Tabs>
        </CardContent>
      </Card>
    </>
  );
}
