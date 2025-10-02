package settlementdetails

import (
	"context"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) CreateSettleDetail(ctx context.Context, tx *gorm.DB, entity *SettlementDetails) error {
	result := tx.WithContext(ctx).Select(
		"settlement_id",
		"transaction_id",
		"transaction_type",
		"procode",
		"mid",
		"tid",
		"card_type",
		"pan",
		"pan_enc",
		"emv_tag",
		"amount",
		"transaction_date",
		"stan",
		"stan_issuer",
		"rrn",
		"trace",
		"batch",
		"trans_mode",
		"bank_code",
		"DE43",
		"response_code",
		"response_at",
		"approval_code",
		"reff_id",
		"DE32",
		"DE33",
		"DE123",
		"issuer_id",
		"status",
		"signature",
		"void_id",
		"cut_off",
		"created_at",
	).Create(&entity)

	return result.Error
}
