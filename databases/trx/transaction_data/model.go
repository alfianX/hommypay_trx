package transactiondata

import "time"

type TransactionData struct {
	ID              int64  		`json:"id"`
	TransactionID   string 		`json:"transaction_id"`
	TransactionType string 		`json:"transaction_type"`
	DataRequest     string 		`json:"data_request"`
	DataResponse    string 		`json:"data_response"`
	IssuerID		int64		`json:"issuer_id"`
	Longitude		string		`json:"longitude"`
	Latitude		string		`json:"latitude"`
	CreatedAt       time.Time 	`gorm:"autoCreateTime:false" json:"created_at"`
	Flag			int64		`json:"flag"`
}

func (TransactionData) TableName() string {
	return "transaction_data"
}