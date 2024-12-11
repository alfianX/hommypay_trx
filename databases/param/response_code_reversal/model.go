package responsecodereversal

type ResponseCodeReversal struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (ResponseCodeReversal) TableName() string {
	return "response_code_reversal"
}