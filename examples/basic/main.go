package main

import (
	"io"
	"log"
	"os"
	"time"

	opensignreport "github.com/TychyInc/opensign-report-pdf"
)

func main() {
	// Basic example with minimum configuration
	config := opensignreport.Config{
		CustomerName:  "テスト株式会社",
		ReceiptNumber: "INV-2024-001",
		IssueDate:     time.Now(),
		MonthlyCount:  50, // 基本料金のみ
		ServiceYear:   2024,
		ServiceMonth:  1,
	}

	reader, err := opensignreport.GenerateInvoice(config)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	file, err := os.Create("INV-2024-001.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("請求書を生成しました: INV-2024-001.pdf")

	// Example with usage fee
	config2 := opensignreport.Config{
		CustomerName:  "サンプル企業",
		ReceiptNumber: "INV-2024-002",
		IssueDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local),
		MonthlyCount:  250, // 基本料金 + 従量課金
		ServiceYear:   2023,
		ServiceMonth:  12,
	}

	reader2, err := opensignreport.GenerateInvoice(config2)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	file2, err := os.Create("INV-2024-002.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()

	_, err = io.Copy(file2, reader2)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("請求書を生成しました: INV-2024-002.pdf")

	// Example with custom fee configuration
	config3 := opensignreport.Config{
		CustomerName:   "カスタム料金設定例",
		ReceiptNumber:  "INV-2024-003",
		IssueDate:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
		MonthlyCount:   150,
		BasicFee:       3000, // 基本料金: 3000円
		FreeUsageCount: 50,   // 無料利用枠: 50件
		UsageFeeRate:   200,  // 従量課金単価: 200円/件
		ServiceYear:    2024,
		ServiceMonth:   1,
	}

	reader3, err := opensignreport.GenerateInvoice(config3)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	file3, err := os.Create("INV-2024-003.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file3.Close()

	_, err = io.Copy(file3, reader3)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("カスタム料金設定の請求書を生成しました: INV-2024-003.pdf")
}
