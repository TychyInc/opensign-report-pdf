package opensignreport

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	gr "github.com/mikeshimura/goreport"
)

// Config holds the configuration for invoice generation
type Config struct {
	CustomerName  string
	ReceiptNumber string
	IssueDate     time.Time
	MonthlyCount  int
	BasicFee      int
	FreeUsageCount int
	UsageFeeRate  int
}

// GenerateInvoice generates an invoice PDF with the given configuration and returns it as io.Reader
func GenerateInvoice(config Config) (io.Reader, error) {
	if config.CustomerName == "" || config.ReceiptNumber == "" || config.MonthlyCount == 0 {
		return nil, fmt.Errorf("customer name, receipt number, and monthly count are required")
	}

	if config.IssueDate.IsZero() {
		config.IssueDate = time.Now()
	}

	basicFee, usageFee, tax, total := calculateFees(config.MonthlyCount, config.BasicFee, config.FreeUsageCount, config.UsageFeeRate)

	// Set default values if not provided
	freeUsageCount := config.FreeUsageCount
	if freeUsageCount == 0 {
		freeUsageCount = 100
	}
	usageFeeRate := config.UsageFeeRate
	if usageFeeRate == 0 {
		usageFeeRate = 100
	}

	data := &ReceiptData{
		CustomerName:  config.CustomerName,
		ReceiptNumber: config.ReceiptNumber,
		IssueDate:     config.IssueDate,
		MonthlyCount:  config.MonthlyCount,
		BasicFee:      basicFee,
		UsageFee:      usageFee,
		Tax:           tax,
		TotalAmount:   total,
		FreeUsageCount: freeUsageCount,
		UsageFeeRate:  usageFeeRate,
	}

	reader, err := generatePDF(data)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func generatePDF(data *ReceiptData) (io.Reader, error) {
	r := gr.CreateGoReport()
	r.PageTotal = true
	r.SumWork["mbreakdown"] = 0.0

	// Set up fonts
	font1 := gr.FontMap{
		FontName: "NotoSansJP",
		FileName: "../../ttf/NotoSansJP-VariableFont_wght.ttf",
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