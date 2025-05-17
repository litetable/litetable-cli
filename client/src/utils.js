/**
 * Unwraps the input data, decodes the value, and maintains the original structure.
 * Ensures single rows include a `rowKey` key.
 * @param {Object} data - The input data object (single row or filter query).
 * @returns {Object} - An object with the same structure as the input, but with decoded values.
 */
export function unwrapAndDecodeData(data) {
  const result = {};

  const processCols = (cols) => {
    const decodedCols = {};
    Object.keys(cols)
      .sort()
      .forEach((family) => {
        decodedCols[family] = {};
        const familyData = cols[family];

        Object.keys(familyData)
          .sort()
          .forEach((qualifier) => {
            decodedCols[family][qualifier] = familyData[qualifier].map(
              (item) => ({
                ...item,
                value: processValue(item.value),
              }),
            );
          });
      });
    return decodedCols;
  };

  // Helper function to process values - first URL-decode, then handle potential JSON
  const processValue = (value) => {
    if (!value) return "";

    try {
      const binary = atob(value);
      const utf8 = new TextDecoder("utf-8").decode(
        new Uint8Array([...binary].map((c) => c.charCodeAt(0))),
      );
      return decodeURIComponent(utf8.replace(/\+/g, " "));
    } catch (e) {
      return value;
    }
  };

  // Process multiple rows
  Object.keys(data)
    .sort()
    .forEach((key) => {
      result[key] = {
        key,
        cols: processCols(data[key].cols),
      };
    });

  return result;
}

export function chunkTextToLiteTableRows({ name, text }, options = {}) {
  const {
    chunkKB = 256,              // Max size of each chunk in kilobytes
    family = "chunks",         // Column family
    qualifier = "contents",        // Qualifier name for chunked text
    prefix = "pdf",            // Row key prefix
  } = options;

  const encoder = new TextEncoder();
  const maxBytes = chunkKB * 1024;

  const chunks = [];
  let currentChunk = '';
  let currentSize = 0;

  const fileBase = name.replace(/\.[^/.]+$/, "")   // remove extension
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")                   // replace non-alphanumerics
    .replace(/(^-|-$)/g, "")

  for (const char of text) {
    const charBytes = encoder.encode(char);
    if (currentSize + charBytes.length > maxBytes) {
      chunks.push(currentChunk);
      currentChunk = '';
      currentSize = 0;
    }

    currentChunk += char;
    currentSize += charBytes.length;
  }

  // Push the final chunk
  if (currentChunk) {
    chunks.push(currentChunk);
  }

  const bunch = chunks.map((chunk, index) => ({
    rowKey: `${prefix}:${fileBase}:${index}`,
    family,
    qualifiers: {
      [qualifier]: chunk,
    },
  }));

  console.log(bunch)
  // Map chunks to LiteTable row format
  return bunch;
}

export function groupItemsByLine(items) {
  const lines = new Map();

  for (const item of items) {
    const y = Math.round(item.transform[5]); // Y position
    if (!lines.has(y)) {
      lines.set(y, []);
    }
    lines.get(y).push(item.str);
  }

  // Sort lines by Y descending (top to bottom on page)
  return [...lines.entries()]
    .sort((a, b) => b[0] - a[0])
    .map(([_, words]) => words.join(" "));
}