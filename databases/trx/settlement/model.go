package settlement

import "time"

type Settlement struct {
	ID               int64     `json:"id"`
	Mid              string    `json:"mid"`
	Tid              string    `json:"tid"`
	Stan             string    `json:"stan"`
	Trace            string    `json:"trace"`
	Batch            string    `json:"batch"`
	RefNo			 string	   `json:"ref_no"`
	CurrencyCode     string    `json:"currency_code"`
	BankID           int64     `json:"bank_id"`
	FirstTrxTime     time.Time `json:"first_trx_time"`
	SubBatchNo       string    `json:"sub_batch_no"`
	SettleDate       string    `json:"settle_date"`
	TotalTransaction int64     `json:"total_transaction"`
	TotalAmount      int64     `json:"total_amount"`
	HostSaleCount    int64     `json:"host_sale_count"`
	HostSaleAmount   int64     `json:"host_sale_amount"`
	HostRefundCount  int64     `json:"host_refund_count"`
	HostRefungAmount int64     `json:"host_refund_amount"`
	PosSaleCount     int64     `json:"pos_sale_count"`
	PosSaleAmount    int64     `json:"pos_sale_amount"`
	PosRefundCount   int64     `json:"pos_refund_count"`
	PosRefundAmount  int64     `json:"pos_refund_amount"`
	Signature		 string	   `json:"signature"`
	ClearingFlag     int64     `json:"clearing_flag"`
	CreatedAt        time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (Settlement) TableName() string {
	return "settlement"
}