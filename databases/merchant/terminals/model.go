package terminals

type Terminals struct {
	ID         int64  `json:"id"`
	MerchantID int64  `json:"merchant_id"`
	Status     string `json:"status"`
	Email      string `json:"email"`
	Batch      string `json:"batch"`
}