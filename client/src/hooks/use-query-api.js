import { useState } from "react";
import { toast } from "sonner";
import { unwrapAndDecodeData } from "@/utils.js";

export function useQueryApi() {
  const [results, setResults] = useState({});
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (payload) => {
    const MIN_SPINNER_DURATION = 250;
    const start = Date.now();
    setIsLoading(true);

    try {
      const response = await fetch("/query", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`Request failed with status: ${response.status}`);
      }

      // Handle CREATE with no decode
      if (payload.type === "CREATE") {
        await delayRemaining(start, MIN_SPINNER_DURATION);
        setIsLoading(false);
        return { success: true };
      }

      const body = await response.json();
      const transformedData = unwrapAndDecodeData(body);
      setResults(transformedData);

      await delayRemaining(start, MIN_SPINNER_DURATION);
      setIsLoading(false);

      toast(`${payload.type} successful`);
      return { success: true, data: transformedData };
    } catch (error) {
      await delayRemaining(start, MIN_SPINNER_DURATION);
      setIsLoading(false);
      toast.error(`${payload.type || "Operation"} failed: ${error.message}`);
      return { success: false, error };
    }
  };

  return {
    handleSubmit,
    results,
    isLoading,
    clearResults: () => setResults({}),
  };
}

const delayRemaining = async (start, minDelay) => {
  const elapsed = Date.now() - start;
  const remaining = Math.max(minDelay - elapsed, 0);
  if (remaining > 0) {
    await new Promise((res) => setTimeout(res, remaining));
  }
};
