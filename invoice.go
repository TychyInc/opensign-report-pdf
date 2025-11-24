package opensignreport

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	gr "github.com/mikeshimura/goreport"
)

// Config holds the configuration for invoice generation
type Config struct {
	CustomerName   string
	ReceiptNumber  string
	IssueDate      time.Time
	MonthlyCount   int
	BasicFee       int
	FreeUsageCount int
	UsageFeeRate   int
	ServiceYear    int
	ServiceMonth   int
}

// GenerateInvoice generates an invoice PDF with the given configuration and returns it as io.Reader
func GenerateInvoice(config Config) (io.Reader, error) {
	if config.CustomerName == "" || config.ReceiptNumber == "" {
		return nil, fmt.Errorf("customer name and receipt number are required")
	}

	if config.IssueDate.IsZero() {
		config.IssueDate = time.Now()
	}

	basicFee, usageFee, total := calculateFees(config.MonthlyCount, config.BasicFee, config.FreeUsageCount, config.UsageFeeRate)

	// Set default values if not provided
	freeUsageCount := config.FreeUsageCount
	if freeUsageCount == 0 {
		freeUsageCount = 100
	}
	usageFeeRate := config.UsageFeeRate
	if usageFeeRate == 0 {
		usageFeeRate = 110 // 税込み
	}

	data := &ReceiptData{
		CustomerName:   config.CustomerName,
		ReceiptNumber:  config.ReceiptNumber,
		IssueDate:      config.IssueDate,
		MonthlyCount:   config.MonthlyCount,
		BasicFee:       basicFee,
		UsageFee:       usageFee,
		TaxRate:        10, // 10%固定
		TotalAmount:    total,
		FreeUsageCount: freeUsageCount,
		UsageFeeRate:   usageFeeRate,
		ServiceYear:    config.ServiceYear,
		ServiceMonth:   config.ServiceMonth,
	}

	reader, err := generatePDF(data)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// getFontPath returns the font file path based on OPENSIGN_TTF_DIR environment variable
// If OPENSIGN_TTF_DIR is set, uses that directory. Otherwise, uses relative path.
func getFontPath() string {
	const fontFileName = "NotoSansJP-VariableFont_wght.ttf"

	if ttfDir := os.Getenv("OPENSIGN_TTF_DIR"); ttfDir != "" {
		return filepath.Join(ttfDir, fontFileName)
	}

	// Default: relative path for backward compatibility
	return filepath.Join("../../ttf", fontFileName)
}

func generatePDF(data *ReceiptData) (io.Reader, error) {
	r := gr.CreateGoReport()
	r.PageTotal = true
	r.SumWork["mbreakdown"] = 0.0

	// Set up fonts
	fontPath := getFontPath()
	font1 := gr.FontMap{
		FontName: "NotoSansJP",
		FileName: fontPath,
	}
	fonts := []*gr.FontMap{&font1}
	r.SetFonts(fonts)

	// Register bands
	r.RegisterBand(gr.Band(ReceiptHeader{data: data}), gr.PageHeader)
	r.RegisterBand(gr.Band(ReceiptDetail{data: data}), gr.Detail)
	r.RegisterBand(gr.Band(ReceiptSummary{data: data}), gr.Summary)
	r.RegisterBand(gr.Band(ReceiptFooter{}), gr.PageFooter)

	// Set up breakdown items as records
	records := []interface{}{
		[]string{fmt.Sprintf("基本利用料 (%d件まで無料)", data.FreeUsageCount), strconv.Itoa(data.BasicFee)},
	}

	if data.UsageFee > 0 {
		billableCount := data.MonthlyCount - data.FreeUsageCount
		records = append(records, []string{
			fmt.Sprintf("従量課金 (%d件 × %d円)", billableCount, data.UsageFeeRate),
			strconv.Itoa(data.UsageFee),
		})
	}

	r.Records = records

	// Configure page
	r.SetPage("A4", "mm", "P")
	r.SetFooterY(280)

	// Generate PDF
	r.Execute("")
	pdfBytes := r.GetBytesPdf()
	return bytes.NewReader(pdfBytes), nil
}
