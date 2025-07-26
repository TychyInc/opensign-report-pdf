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
}

func calculateFees(monthlyCount int) (basicFee, usageFee, tax, total int) {
	basicFee = 2000

	if monthlyCount > 100 {
		usageFee = (monthlyCount - 100) * 100
	}

	subtotal := basicFee + usageFee
	tax = subtotal * 10 / 100
	total = subtotal + tax

	return basicFee, usageFee, tax, total
}
