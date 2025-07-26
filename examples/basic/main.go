package main

import (
	"log"
	"time"

	opensignreport "github.com/ryuyama/opensign-report-pdf"
)

func main() {
	// Basic example with minimum configuration
	config := opensignreport.Config{
		CustomerName:  "テスト株式会社",
		ReceiptNumber: "INV-2024-001",
		IssueDate:     time.Now(),
		MonthlyCount:  50, // 基本料金のみ
	}

	filename, err := opensignreport.GenerateInvoice(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("請求書を生成しました: %s", filename)

	// Example with usage fee
	config2 := opensignreport.Config{
		CustomerName:  "サンプル企業",
		ReceiptNumber: "INV-2024-002",
		IssueDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
		MonthlyCount:  250, // 基本料金 + 従量課金
	}

	filename2, err := opensignreport.GenerateInvoice(config2)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("請求書を生成しました: %s", filename2)
}