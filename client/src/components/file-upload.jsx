import { useState, useCallback } from "react";
import { pdfjs } from 'react-pdf';
import { Upload, X, Check, Loader } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const supportedFileTypes = ["application/pdf"];
const maxKiloBytes = 10240;

pdfjs.GlobalWorkerOptions.workerSrc = new URL(
	'pdfjs-dist/build/pdf.worker.min.mjs',
	import.meta.url,
).toString();

export function FileUploadInput({
  acceptedFileTypes = supportedFileTypes,
  maxSize = maxKiloBytes,
  file,
  setFile,
	cb = null,
}) {
  const [dragActive, setDragActive] = useState(false);
  const [validationMessage, setValidationMessage] = useState("");

  const formatSize = (size) =>
    size > 1000 ? `${Math.round(size / 1024)} MB` : `${size} KB`;

  const onDragOver = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(true);
  }, []);

  const onDragLeave = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
  }, []);

  const onDrop = useCallback(
    (e) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);

      const droppedFile = e.dataTransfer.files[0];
      if (droppedFile) {
        if (acceptedFileTypes.includes(droppedFile.type)) {
          if (droppedFile.size / 1024 > maxSize) {
            setValidationMessage(
              `File size exceeds the maximum limit of ${formatSize(maxSize)}.`,
            );
            return;
          }
          setFile(droppedFile);

          const reader = new FileReader();
          reader.onload = (event) => {
            const arrayBuffer = event.target.result;
            const bytes = new Uint8Array(arrayBuffer);
            setFile({
              name: droppedFile.name,
              type: droppedFile.type,
              size: droppedFile.size,
              bytes,
            });
          };
          reader.readAsArrayBuffer(droppedFile);
          setValidationMessage("");
        } else {
          setValidationMessage(
            "Unsupported file type. Please upload a PDF file.",
          );
        }
      }
    },
    [acceptedFileTypes, maxSize],
  );

  const onFileInputChange = useCallback(
    (e) => {
      if (e.target.files && e.target.files[0]) {
        const selectedFile = e.target.files[0];
        if (acceptedFileTypes.includes(selectedFile.type)) {
          if (selectedFile.size / 1024 > maxSize) {
            setValidationMessage(
              `File size exceeds the maximum limit of ${formatSize(maxSize)}.`,
            );
            return;
          }
          setFile(selectedFile);

          const reader = new FileReader();
          reader.onload = async (event) => {
            const arrayBuffer = event.target.result;
            const bytes = new Uint8Array(arrayBuffer);
            setFile({
              name: selectedFile.name,
              type: selectedFile.type,
              size: selectedFile.size,
              bytes,
            });
          };
          reader.readAsArrayBuffer(selectedFile);
          setValidationMessage("");
        } else {
          setValidationMessage(
            "Unsupported file type. Please upload a PDF file.",
          );
        }
      }
    },
    [acceptedFileTypes, maxSize],
  );

  const removeFile = useCallback(() => {
    setFile(null);
    setFile(null);
    setValidationMessage("");
		if (cb) {
			cb();
		}
  }, []);

  return (
    <form className="max-w-md">
      <div
        className={`relative border-2 border-dashed rounded-lg p-6 ${
          dragActive ? "border-primary" : "border-muted-foreground"
        }`}
        onDragOver={onDragOver}
        onDragLeave={onDragLeave}
        onDrop={onDrop}
      >
        <Input
          id="file-upload"
          type="file"
          accept={acceptedFileTypes.join(",")}
          className="sr-only"
          onChange={onFileInputChange}
          aria-label="File upload"
          disabled={!!file}
        />
        <Label
          htmlFor="file-upload"
          className={`flex flex-col items-center justify-center cursor-pointer ${file ? "cursor-not-allowed" : ""}`}
        >
          <Upload className="w-10 h-10 text-muted-foreground mb-2" />
          <p className="text-sm text-muted-foreground mb-2">
            Drag & drop a single file here or just click to select
          </p>
          <Button
            type="button"
            variant="outline"
            onClick={() => document.getElementById("file-upload")?.click()}
            size="sm"
            disabled={!!file}
          >
            Select File
          </Button>
        </Label>
      </div>
      {validationMessage && (
        <p className="text-red-500 text-sm mt-2">{validationMessage}</p>
      )}
      {file && (
        <div className="mt-4 flex items-center">
          <Check className="w-5 h-5 text-green-600 font-bold mr-2" />
          <div className="flex items-center justify-between text-sm bg-muted p-2 rounded w-full">
            <span className="truncate">{file.name}</span>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              onClick={removeFile}
              aria-label={`Remove ${file.name}`}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
        </div>
      )}
    </form>
  );
}

export default function FileUploadForm({ fileCB = null }) {
	const [file, setFile] = useState(null);
	const [uploaded, setUploaded] = useState(false);
	const [isProcessed, setIsProcessed] = useState(false);

	const handleSubmit = async (e) => {
		e.preventDefault();
		setUploaded(true);

		if (file && file.bytes) {
			const pdfDocument = await pdfjs.getDocument(new Uint8Array(file.bytes)).promise;
			const pages = pdfDocument.numPages;

			const allTextItems = await Promise.all(
				Array.from({ length: pages }, (_, i) =>
					pdfDocument.getPage(i + 1).then((page) => page.getTextContent())
				)
			);

			const textItems = allTextItems.flatMap((textContent) => textContent.items);
			const contents = textItems.map((item) => item.str);
			const text = contents.join("");

			const updatedFile = { ...file, text };
			setFile(updatedFile);

			// ✅ Call the callback AFTER setFile, allow it to be async
			if (fileCB) {
				try {
					await fileCB(updatedFile); // Allow the parent to do async work
				} catch (err) {
					console.error("File callback failed:", err);
				}
			}

			setIsProcessed(true);
		}

		setUploaded(false);
	};

	const reset = () => {
			setFile(null);
			setUploaded(false);
			setIsProcessed(false);
		}
	return (
		<>
			<FileUploadInput file={file} setFile={setFile} cb={reset} />
			<Button className="w-full mt-4" disabled={!file || isProcessed} onClick={handleSubmit}>
				{/* eslint-disable-next-line no-nested-ternary */}
				{isProcessed ? (
					'Complete'
				) : uploaded ? (
					<Loader className="animate-spin" size={16} />
				) : (
					'Upload file'
				)}
			</Button>
		</>
	);
}