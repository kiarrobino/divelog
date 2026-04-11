package calculator

// ndlTable maps depth in metres to no-decompression limit in minutes.
// Values sourced from PADI Recreational Dive Planner tables.
var ndlTable = map[float64]int{
	10: 219,
	12: 147,
	14: 98,
	16: 72,
	18: 56,
	20: 45,
	25: 29,
	30: 20,
	35: 14,
	40: 9,
}

// Calculate returns the no-decompression limit in minutes for a given depth.
// Depth is rounded to the nearest entry in the NDL table. If the depth
// exceeds 40m, 0 is returned — recreational diving does not go beyond 40m.
func Calculate(depth float64) int {
	// Round up to the nearest table depth — always conservative.
	for _, tableDepth := range []float64{10, 12, 14, 16, 18, 20, 25, 30, 35, 40} {
		if depth <= tableDepth {
			return ndlTable[tableDepth]
		}
	}
	return 0
}
