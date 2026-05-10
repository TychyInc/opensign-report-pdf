package opensignreport

import "testing"

func TestCalculateFees(t *testing.T) {
	tests := []struct {
		name           string
		monthlyCount   int
		basicFeeParam  int
		freeUsageCount int
		usageFeeRate   int
		wantBasicFee   int
		wantUsageFee   int
		wantTotal      int
	}{
		{
			name:           "monthly count is zero",
			monthlyCount:   0,
			basicFeeParam:  0,
			freeUsageCount: 0,
			usageFeeRate:   0,
			wantBasicFee:   2200,
			wantUsageFee:   0,
			wantTotal:      2200,
		},
		{
			name:           "monthly count within free usage",
			monthlyCount:   50,
			basicFeeParam:  0,
			freeUsageCount: 100,
			usageFeeRate:   110,
			wantBasicFee:   2200,
			wantUsageFee:   0,
			wantTotal:      2200,
		},
		{
			name:           "monthly count exceeds free usage",
			monthlyCount:   150,
			basicFeeParam:  0,
			freeUsageCount: 100,
			usageFeeRate:   110,
			wantBasicFee:   2200,
			wantUsageFee:   5500,
			wantTotal:      7700,
		},
		{
			name:           "custom basic fee with zero monthly count",
			monthlyCount:   0,
			basicFeeParam:  3300,
			freeUsageCount: 0,
			usageFeeRate:   0,
			wantBasicFee:   3300,
			wantUsageFee:   0,
			wantTotal:      3300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBasicFee, gotUsageFee, gotTotal := calculateFees(
				tt.monthlyCount,
				tt.basicFeeParam,
				tt.freeUsageCount,
				tt.usageFeeRate,
			)
			if gotBasicFee != tt.wantBasicFee {
				t.Errorf("basicFee = %d, want %d", gotBasicFee, tt.wantBasicFee)
			}
			if gotUsageFee != tt.wantUsageFee {
				t.Errorf("usageFee = %d, want %d", gotUsageFee, tt.wantUsageFee)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("total = %d, want %d", gotTotal, tt.wantTotal)
			}
		})
	}
}
