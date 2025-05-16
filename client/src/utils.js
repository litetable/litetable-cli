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
    chunkSize = 500,         // Number of characters per chunk
    family = "chunks",
    qualifier = "text",
    prefix = "pdf",
  } = options;

  // Sanitize and slugify filename
  const fileBase = name.replace(/\.[^/.]+$/, "")   // remove extension
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")                   // replace non-alphanumerics
    .replace(/(^-|-$)/g, "")                       // trim hyphens

  const rowPrefix = `${prefix}:${fileBase}`

  // Split text into chunks
  const chunks = []
  for (let i = 0; i < text.length; i += chunkSize) {
    const chunkText = text.slice(i, i + chunkSize).trim()
    if (!chunkText) continue

    chunks.push({
      rowKey: `${rowPrefix}:${chunks.length}`,  // maintains order
      family,
      qualifiers: {
        [qualifier]: chunkText
      }
    })
  }

  return chunks
}
