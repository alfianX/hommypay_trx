package types

type SaleRequest struct {
	PaymentInformation struct {
		Procode         string `json:"procode" binding:"required,min=6,max=6,numeric"`
		TID             string `json:"tid" binding:"required,min=8,max=8"`
		MID             string `json:"mid" binding:"required,min=15,max=15"`
		Amount          int64  `json:"amount" binding:"required,numeric"`
		Tip             int64  `json:"tip" binding:"numeric"`
		STAN            string `json:"stan" binding:"required,min=6,max=6,numeric"`
		Trace           string `json:"trace" binding:"required,min=6,max=6,numeric"`
		Batch           string `json:"batch" binding:"required,min=6,max=6,numeric"`
		TransactionDate string `json:"transactionDate" binding:"required"`
		KSN             string `json:"ksn"`
	} `json:"paymentInformation" binding:"required"`
	CardInformation struct {
		AID        string `json:"aid"`
		PAN        string `json:"pan" binding:"required,min=16,max=19,numeric"`
		Expiry     string `json:"expiry"`
		TrackData2 string `json:"trackData" binding:"required"`
		EMVTag     string `json:"emvTag"`
		PinBlock   string `json:"pinBlock"`
	} `json:"cardInformation" binding:"required"`
	PosTerminal struct {
		TransMode string `json:"transMode" binding:"required,min=3,numeric"`
		Code      string `json:"code"`
		KeyMode   int    `json:"keyMode"`
	} `json:"posTerminal" binding:"required"`
	ISO8583 string `json:"ISO8583"`
}

type SaleLogRequest struct {
	PaymentInformation struct {
		Procode         string `json:"procode" binding:"required,min=6,max=6,numeric"`
		TID             string `json:"tid" binding:"required,min=8,max=8"`
		MID             string `json:"mid" binding:"required,min=15,max=15"`
		Amount          int64  `json:"amount" binding:"required,numeric"`
		Tip             int64  `json:"tip" binding:"numeric"`
		STAN            string `json:"stan" binding:"required,min=6,max=6,numeric"`
		Trace           string `json:"trace" binding:"required,min=6,max=6,numeric"`
		Batch           string `json:"batch" binding:"required,min=6,max=6,numeric"`
		TransactionDate string `json:"transactionDate" binding:"required"`
		KSN             string `json:"ksn"`
	} `json:"paymentInformation" binding:"required"`
	PosTerminal struct {
		TransMode string `json:"transMode" binding:"required,min=3,numeric"`
		Code      string `json:"code"`
		KeyMode   int    `json:"keyMode"`
	} `json:"posTerminal" binding:"required"`
}