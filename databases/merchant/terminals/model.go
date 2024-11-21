package terminals

type Terminals struct {
	TerminalID int64  `json:"terminal_id"`
	MerchantID int64  `json:"merchant_id"`
	Status     string `json:"status"`
	Email      string `json:"email"`
}