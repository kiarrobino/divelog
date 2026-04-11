package calculator

import "testing"

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		depth   float64
		wantNDL int
	}{
		{"10m", 10, 219},
		{"12m", 12, 147},
		{"14m", 14, 98},
		{"16m", 16, 72},
		{"18m", 18, 56},
		{"20m", 20, 45},
		{"25m", 25, 29},
		{"30m", 30, 20},
		{"35m", 35, 14},
		{"40m", 40, 9},
		{"0m", 0, 219},
		{"45m", 45, 0},
		{"32m", 32, 14},
		{"-5m", -5, 219},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.depth)
			if got != tt.wantNDL {
				t.Errorf("Calculate(%v) = %v, want %v", tt.depth, got, tt.wantNDL)
			}
		})
	}
}
