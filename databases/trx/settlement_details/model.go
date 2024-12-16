package settlementdetails

import "time"

type SettlementDetails struct {
	ID                  int64     `json:"id"`
	SettlementID        int64     `json:"settlement_id"`
	TransactionID       string    `json:"transaction_id"`
	TransactionType     string    `json:"transaction_type"`
	Procode				string	  `json:"procode"`
	Mid                 string    `json:"mid"`
	Tid                 string    `json:"tid"`
	CardType            string    `json:"card_type"`
	Pan                 string    `json:"pan"`
	PanEnc				string	  `json:"pan_enc"`
	TrackData           string    `json:"track_data"`
	EmvTag              string    `json:"emv_tag"`
	Amount              int64     `json:"amount"`
	TransactionDate     time.Time `gorm:"autoCreateTime:false" json:"transaction_date"`
	Stan                string    `json:"stan"`
	StanIssuer			string	  `json:"stan_issuer"`
	Rrn					string	  `json:"rrn"`
	Trace               string    `json:"trace"`
	Batch				string	  `json:"batch"`
	TransMode			string	  `json:"trans_mode"`
	IsoRequest      	string    `json:"iso_request"`
	IsoRequestIssuer  	string    `json:"iso_request_issuer"`
	ResponseCode        string    `json:"response_code"`
	ResponseAt          time.Time `gorm:"autoCreateTime:false" json:"response_at"`
	ApprovalCode        string    `json:"approval_code"`
	RefID               string    `json:"reff_id"`
	DE32				string	  `json:"DE32"`
	DE33				string	  `json:"DE33"`
	DE123				string	  `json:"DE123"`
	IsoResponse     	string    `json:"iso_response"`
	IsoResponseIssuer 	string    `json:"iso_response_issuer"`
	IssuerID            int64     `json:"issuer_id"`
	Status              int64     `json:"status"`
	Signature			string	  `json:"signature"`
	VoidID              string    `json:"void_id"`
	CutOff				string	  `json:"cut_off"`
	CreatedAt           time.Time `gorm:"autoCreateTime:false" json:"created_at"`
}