package server

import (
	"github.com/litetable/litetable-cli/internal/litetable"
	"github.com/litetable/litetable-db/pkg/proto"
)

func unwrap(rows map[string]*proto.Row) map[string]*litetable.Row {
	result := make(map[string]*litetable.Row)

	for key, row := range rows {
		//
		ltRow := &litetable.Row{
			Key:     key,
			Columns: make(map[string]litetable.VersionedQualifier),
		}
		cols := row.GetCols()

		// for each column family
		for family, qualifierData := range cols {
			// Initialize the family in our result if not exists
			if _, exists := ltRow.Columns[family]; !exists {
				ltRow.Columns[family] = make(litetable.VersionedQualifier)
			}

			// For each qualifier in this family
			for qualifier, values := range qualifierData.GetQualifiers() {
				tsValues := make([]litetable.TimestampedValue, 0, len(values.GetValues()))

				// For each value in this qualifier
				for _, val := range values.GetValues() {
					tsValues = append(tsValues, litetable.TimestampedValue{
						Value:     val.GetValue(),
						Timestamp: val.GetTimestampUnix(),
					})
				}

				ltRow.Columns[family][qualifier] = tsValues
			}
		}

		result[key] = ltRow
	}

	return result
}
