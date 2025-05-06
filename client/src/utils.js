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
                value: atob(item.value), // Decode Base64
              }),
            );
          });
      });
    return decodedCols;
  };

  // Multiple rows case
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
