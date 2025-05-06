/**
 * Transforms the input data into a format suitable for a shadcn table.
 * @param {Object} data - The input data object (single row or filter query).
 * @returns {Object} - An object containing table rows and metadata.
 */
export function transformToTableData(data) {
	const tableData = [];
	const columnFamilies = new Set();
	const qualifiers = new Set();

	// Normalize single row response to match filter query structure
	const rows = data.key ? [data] : Object.values(data);

	rows.forEach((row) => {
		const { key, cols } = row;

		for (const family in cols) {
			if (Object.prototype.hasOwnProperty.call(cols, family)) {
				columnFamilies.add(family);
				const familyData = cols[family];

				for (const qualifier in familyData) {
					if (Object.prototype.hasOwnProperty.call(familyData, qualifier)) {
						qualifiers.add(qualifier);
						const qualifierData = familyData[qualifier];

						qualifierData.forEach((item) => {
							tableData.push({
								key,
								family,
								qualifier,
								value: atob(item.value), // Decode Base64
								timestamp: item.timestamp,
							});
						});
					}
				}
			}
		}
	});

	return {
		tableData,
		metadata: {
			columnFamilies: Array.from(columnFamilies),
			qualifiers: Array.from(qualifiers),
		},
	};
}

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
		Object.keys(cols).sort().forEach((family) => {
			decodedCols[family] = {};
			const familyData = cols[family];

			Object.keys(familyData).sort().forEach((qualifier) => {
				decodedCols[family][qualifier] = familyData[qualifier].map((item) => ({
					...item,
					value: atob(item.value), // Decode Base64
				}));
			});
		});
		return decodedCols;
	};

	if (data.key) {
		// Single row case
		result[data.key] = {
			key: data.key,
			cols: processCols(data.cols),
		};
	} else {
		// Multiple rows case
		Object.keys(data).sort().forEach((key) => {
			result[key] = {
				key,
				cols: processCols(data[key].cols),
			};
		});
	}

	return result;
}