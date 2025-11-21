package opensignreport

import (
	"time"
)

type ReceiptData struct {
	CustomerName   string
	ReceiptNumber  string
	IssueDate      time.Time
	MonthlyCount   int
	BasicFee       int // 税込み
	UsageFee       int // 税込み
	TaxRate        int // 税率（パーセント）
	TotalAmount    int // 税込み合計
	FreeUsageCount int
	UsageFeeRate   int // 税込み単価
	ServiceYear    int
	ServiceMonth   int
}

func calculateFees(monthlyCount, basicFeeParam, freeUsageCount, usageFeeRate int) (basicFee, usageFee, total int) {
	// デフォルト値: 税込み価格
	if basicFeeParam == 0 {
		basicFee = 2200 // 税込み
	} else {
		basicFee = basicFeeParam
	}

	if freeUsageCount == 0 {
		freeUsageCount = 100
	}

	if usageFeeRate == 0 {
		usageFeeRate = 110 // 税込み
	}

	if monthlyCount > freeUsageCount {
		usageFee = (monthlyCount - freeUsageCount) * usageFeeRate
	}

	total = basicFee + usageFee

	return basicFee, usageFee, total
}
