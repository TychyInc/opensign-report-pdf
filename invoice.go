package opensignreport

import (
	"fmt"
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
}

// GenerateInvoice generates an invoice PDF with the given configuration
func GenerateInvoice(config Config) (string, error) {
	if config.CustomerName == "" || config.ReceiptNumber == "" || config.MonthlyCount == 0 {
		return "", fmt.Errorf("customer name, receipt number, and monthly count are required")
	}

	if config.IssueDate.IsZero() {
		config.IssueDate = time.Now()
	}

	basicFee, usageFee, tax, total := calculateFees(config.MonthlyCount)

	data := &ReceiptData{
		CustomerName:  config.CustomerName,
		ReceiptNumber: config.ReceiptNumber,
		IssueDate:     config.IssueDate,
		MonthlyCount:  config.MonthlyCount,
		BasicFee:      basicFee,
		UsageFee:      usageFee,
		Tax:           tax,
		TotalAmount:   total,
	}

	filename := fmt.Sprintf("%s.pdf", data.ReceiptNumber)
	err := generatePDF(data, filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func generatePDF(data *ReceiptData, filename string) error {
	r := gr.CreateGoReport()
	r.PageTotal = true
	r.SumWork["mbreakdown"] = 0.0

	// Set up fonts
	font1 := gr.FontMap{
		FontName: "NotoSansJP",
		FileName: "ttf//NotoSansJP-VariableFont_wght.ttf",
	}
	fonts := []*gr.FontMap{&font1}
	r.SetFonts(fonts)

	// Register bands
	r.RegisterBand(gr.Band(ReceiptHeader{data: data}), gr.PageHeader)
	r.RegisterBand(gr.Band(ReceiptDetail{data: data}), gr.Detail)
	r.RegisterBand(gr.Band(ReceiptSummary{data: data}), gr.Summary)
	r.RegisterBand(gr.Band(ReceiptFooter{}), gr.PageFooter)

	// Set up breakdown items as records
	r.Records = []interface{}{
		[]string{"基本利用料", strconv.Itoa(data.BasicFee)},
		[]string{fmt.Sprintf("利用料 (%d件 × 100円)", data.MonthlyCount), strconv.Itoa(data.UsageFee)},
	}

	// Configure page
	r.SetPage("A4", "mm", "P")
	r.SetFooterY(280)

	// Generate PDF
	r.Execute(filename)
	return nil
}