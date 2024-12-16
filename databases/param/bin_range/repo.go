package binrange

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) GetUrlByPAN(ctx context.Context, pan, cardType string) (int64, string, int64, string, error) {
	var data BinRange

	result := r.Db.WithContext(ctx).Preload("Issuer").Limit(1).
				Where("(? BETWEEN bin_range.pan_range_low AND bin_range.pan_range_high) AND card_type = ?", 
				pan, cardType).Find(&data)

	
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, "", 0, "",  result.Error
	}
	
	return data.IssuerID, data.Issuer.IssuerName, data.Issuer.IssuerConnType, data.Issuer.IssuerService,  nil
}

func (r Repo) GetCardTypeByPAN(ctx context.Context, pan string) (string, error) {
	var data BinRange

	result := r.Db.WithContext(ctx).Where("? BETWEEN bin_range.pan_range_low AND bin_range.pan_range_high", pan).First(&data)

	return data.CardType, result.Error
}