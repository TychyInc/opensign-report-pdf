package opensignreport

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Set up font path for tests
	os.Setenv("OPENSIGN_TTF_DIR", "./ttf")
	os.Setenv("OPENSIGN_TTF_FILE", "NotoSansJP-VariableFont_wght.ttf")
	os.Exit(m.Run())
}

func TestGenerateInvoice(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		wantErr     bool
		wantErrMsg  string
		wantNonNil  bool
	}{
		{
			name: "monthly count is zero",
			config: Config{
				CustomerName:   "Test Customer",
				ReceiptNumber:  "R-001",
				MonthlyCount:   0,
				IssueDate:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ServiceYear:    2024,
				ServiceMonth:   1,
			},
			wantErr:    false,
			wantNonNil: true,
		},
		{
			name: "monthly count within free usage",
			config: Config{
				CustomerName:   "Test Customer",
				ReceiptNumber:  "R-002",
				MonthlyCount:   50,
				FreeUsageCount: 100,
				IssueDate:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ServiceYear:    2024,
				ServiceMonth:   1,
			},
			wantErr:    false,
			wantNonNil: true,
		},
		{
			name: "monthly count exceeds free usage",
			config: Config{
				CustomerName:   "Test Customer",
				ReceiptNumber:  "R-003",
				MonthlyCount:   150,
				FreeUsageCount: 100,
				UsageFeeRate:   110,
				IssueDate:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ServiceYear:    2024,
				ServiceMonth:   1,
			},
			wantErr:    false,
			wantNonNil: true,
		},
		{
			name: "missing customer name",
			config: Config{
				ReceiptNumber: "R-004",
				MonthlyCount:  50,
			},
			wantErr:    true,
			wantErrMsg: "customer name and receipt number are required",
		},
		{
			name: "missing receipt number",
			config: Config{
				CustomerName: "Test Customer",
				MonthlyCount: 50,
			},
			wantErr:    true,
			wantErrMsg: "customer name and receipt number are required",
		},
		{
			name: "default values applied",
			config: Config{
				CustomerName:  "Test Customer",
				ReceiptNumber: "R-005",
				MonthlyCount:  0,
				IssueDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				ServiceYear:   2024,
				ServiceMonth:  1,
			},
			wantErr:    false,
			wantNonNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := GenerateInvoice(tt.config)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantErrMsg != "" && err.Error() != tt.wantErrMsg {
					t.Errorf("error message = %q, want %q", err.Error(), tt.wantErrMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantNonNil && reader == nil {
				t.Errorf("expected non-nil reader, got nil")
			}

			// Verify we can read from the reader
			if reader != nil {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					t.Errorf("failed to read from reader: %v", err)
				}
				if n == 0 {
					t.Errorf("reader returned 0 bytes")
				}
			}
		})
	}
}
