package reversals

import "time"

type Reversals struct {
	ID              int64  		`json:"id"`
	TransactionID   string 		`json:"transaction_id"`
	TransactionType	string		`json:"transaction_type"`
	Procode			string		`json:"procode"`
	Mid				string		`json:"mid"`
	Tid				string		`json:"tid"`
	Amount			int64		`json:"amount"`
	TransactionDate time.Time	`gorm:"autoCreateTime:false" json:"transaction_date"`
	Stan			string		`json:"stan"`
	Trace			string		`json:"trace"`
	Batch			string		`json:"batch"`
	IsoRequest      string 		`json:"iso_request"`
	ResponseCodeOrg string 		`json:"response_code_origin"`
	ResponseCode    string 		`json:"response_code"`
	IsoResponse     string 		`json:"iso_response"`
	RepeatCount     int64  		`json:"repeat_count"`
	Flag			int64		`json:"flag"`
	CreatedAt       time.Time	`gorm:"autoCreateTime:false" json:"created_at"`
}