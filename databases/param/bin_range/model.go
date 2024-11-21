package binrange

import "github.com/alfianX/hommypay_trx/databases/param/issuer"

type BinRange struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	BankCode     string        `json:"bank_code"`
	CardType     string        `json:"card_type"`
	PanRangeLow  string        `json:"pan_range_low"`
	PanRangeHigh string        `json:"pan_range_high"`
	IssuerID     int64         `json:"issuer_id"`
	Issuer       issuer.Issuer `gorm:"foreignKey:IssuerID"`
}

func (BinRange) TableName() string {
	return "bin_range"
}