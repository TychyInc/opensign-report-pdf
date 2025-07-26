package opensignreport

import (
	"strconv"

	gr "github.com/mikeshimura/goreport"
)

const (
	FONT_SIZE_TITLE  = 16
	FONT_SIZE_NORMAL = 10
	FONT_SIZE_SMALL  = 9
	LINE_HEIGHT      = 5

	// Column positions for table
	col0Left   = 20
	col1Left   = 100
	col2Left   = 170
	tableWidth = 150
)

// ReceiptHeader creates the main header with title and company info
type ReceiptHeader struct {
	data *ReceiptData
}

func (h ReceiptHeader) GetHeight(report gr.GoReport) float64 {
	return 90
}

func (h ReceiptHeader) Execute(report gr.GoReport) {
	// Receipt number and date
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	dateStr := h.data.IssueDate.Format("2006年01月02日")
	report.CellRight(190, 10, 0, dateStr)
	report.CellRight(190, 15, 0, "No. "+h.data.ReceiptNumber)

	// Title
	report.Font("NotoSansJP", FONT_SIZE_TITLE, "")
	report.Cell(85, 25, "請　求　書")

	// Company info (right side)
	report.Font("NotoSansJP", 11, "")
	x := 110.0
	report.Cell(x, 40, "株式会社ティヒ")
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	report.Cell(x, 48, "〒152-0004")
	report.Cell(x, 52, "東京都目黒区鷹番3-21-14 303")
	report.Cell(x, 58, "登録番号: T9011001155105")

	// Customer name (left side)
	report.Font("NotoSansJP", FONT_SIZE_NORMAL+2, "")
	customerName := h.data.CustomerName + " 御中"
	report.Cell(col0Left, 40, customerName)

	// Underline for customer name
	report.LineType("straight", 0.3)
	report.GrayStroke(0.5)
	// Calculate text width for underline
	report.Converter.Pdf.SetFont("NotoSansJP", "", 12)
	w, _ := report.Converter.Pdf.MeasureTextWidth(customerName)
	report.LineH(col0Left, 44, col0Left+w/report.ConvPt)

	// Total amount
	report.GrayStroke(0.0)
	report.Font("NotoSansJP", FONT_SIZE_NORMAL, "")
	report.Cell(col0Left, 55, "下記の通りご請求申し上げます。")

	report.Font("NotoSansJP", 14, "")
	report.Cell(col0Left, 68, "ご請求金額")
	totalStr := "￥" + gr.AddComma(strconv.Itoa(h.data.TotalAmount)) + " -"
	report.CellRight(col0Left+80, 68, 0, totalStr)
	report.LineH(col0Left, 72, col0Left+80)

	// Description
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	report.Cell(col0Left, 80, "但し、OpenSign利用料として")
}

// ReceiptDetail shows the breakdown table
type ReceiptDetail struct {
	data *ReceiptData
}

func (d ReceiptDetail) GetHeight(report gr.GoReport) float64 {
	return 6
}

func (d ReceiptDetail) Execute(report gr.GoReport) {
	cols := report.Records[report.DataPos].([]string)
	y := 1.5
	x := float64(col0Left)

	// Draw table borders
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x, 0, x+tableWidth, 6)
	report.LineV(x, 0, 6)
	report.LineV(x+80, 0, 6)
	report.LineV(x+tableWidth, 0, 6)

	// Content
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	report.Cell(x+2, y, cols[0])
	report.CellRight(x+tableWidth-2, y, 0, "￥"+gr.AddComma(cols[1]))

	// Sum for breakdown
	amount, _ := strconv.ParseFloat(cols[1], 64)
	report.SumWork["mbreakdown"] += amount
}

// ReceiptSummary shows the summary with tax and total
type ReceiptSummary struct {
	data *ReceiptData
}

func (s ReceiptSummary) GetHeight(report gr.GoReport) float64 {
	return 35
}

func (s ReceiptSummary) Execute(report gr.GoReport) {
	x := float64(col0Left)
	y := 0.0

	// Draw summary table
	report.LineType("straight", 0.3)
	report.GrayStroke(0)

	// Subtotal row
	report.Rect(x, y, x+tableWidth, y+6)
	report.LineV(x+80, y, y+6)
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	report.Cell(x+2, y+1.5, "小計")
	subtotal := s.data.BasicFee + s.data.UsageFee
	report.CellRight(x+tableWidth-2, y+1.5, 0, "￥"+gr.AddComma(strconv.Itoa(subtotal)))

	// Tax row
	y += 6
	report.Rect(x, y, x+tableWidth, y+6)
	report.LineV(x+80, y, y+6)
	report.Cell(x+2, y+1.5, "消費税（10%）")
	report.CellRight(x+tableWidth-2, y+1.5, 0, "￥"+gr.AddComma(strconv.Itoa(s.data.Tax)))

	// Total row with emphasis
	y += 6
	report.GrayStroke(0.85)
	report.LineH(x, y, x+tableWidth)
	report.LineType("straight", 0.3)
	report.GrayStroke(0)
	report.Rect(x, y, x+tableWidth, y+6)
	report.LineV(x+80, y, y+6)
	report.Font("NotoSansJP", FONT_SIZE_NORMAL, "")
	report.Cell(x+2, y+1.5, "合計（税込）")
	report.CellRight(x+tableWidth-2, y+1.5, 0, "￥"+gr.AddComma(strconv.Itoa(s.data.TotalAmount)))
}

// ReceiptFooter for page numbers
type ReceiptFooter struct{}

func (f ReceiptFooter) GetHeight(report gr.GoReport) float64 {
	return 10
}

func (f ReceiptFooter) Execute(report gr.GoReport) {
	report.Font("NotoSansJP", FONT_SIZE_SMALL, "")
	report.Cell(95, 5, "Page "+strconv.Itoa(report.Page))
}
