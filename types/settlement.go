package types

type SettlementRequest struct {
	SettlementType     string `json:"settlementType" binding:"required"`
	PaymentInformation struct {
		TID        string `json:"tid" binding:"required"`
		MID        string `json:"mid" binding:"required"`
		STAN       string `json:"stan" binding:"required"`
		Trace      string `json:"trace" binding:"required"`
		Batch      string `json:"batch" binding:"required"`
		SettleDate string `json:"settleDate" binding:"required"`
	} `json:"paymentInformation" binding:"required"`
	OrderInformation struct {
		TotalTransaction int64 `json:"totalTransaction"`
		TotalAmount      int64 `json:"totalAmount" `
		SaleCount        int64 `json:"saleCount" `
		SaleAmount       int64 `json:"saleAmount" `
		VoidCount        int64 `json:"voidCount" `
		VoidAmount       int64 `json:"voidAmount" `
	} `json:"orderInformation" binding:"required"`
}