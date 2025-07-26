package opensignreport

import (
	"time"
)

type ReceiptData struct {
	CustomerName  string
	ReceiptNumber string
	IssueDate     time.Time
	MonthlyCount  int
	BasicFee      int
	UsageFee      int
	Tax           int
	TotalAmount   int
	FreeUsageCount int
	UsageFeeRate  int
	ServiceYear   int
	ServiceMonth  int
}

func calculateFees(monthlyCount, basicFeeParam, freeUsageCount, usageFeeRate int) (basicFee, usageFee, tax, total int) {
	if basicFeeParam == 0 {
		basicFee = 2000
	} else {
		basicFee = basicFeeParam
	}

	if freeUsageCount == 0 {
		freeUsageCount = 100
	}

	if usageFeeRate == 0 {
		usageFeeRate = 100
	}

	if monthlyCount > freeUsageCount {
		usageFee = (monthlyCount - freeUsageCount) * usageFeeRate
	}

	subtotal := basicFee + usageFee
	tax = subtotal * 10 / 100
	total = subtotal + tax

	return basicFee, usageFee, tax, total
}
