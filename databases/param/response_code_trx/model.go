package responsecodetrx

type ResponseCodeTrx struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (ResponseCodeTrx) TableName() string {
	return "response_code_trx"
}