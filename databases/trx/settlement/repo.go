package settlement

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) CreateSettle(ctx context.Context, tx *gorm.DB, entity *Settlement) error {
	result := tx.WithContext(ctx).Select(
		"settlement_id",
		"mid",
		"tid",
		"stan",
		"trace",
		"batch",
		"ref_no",
		"sub_batch_no",
		"settle_date",
		"total_transaction",
		"total_amount",
		"host_sale_count",
		"host_sale_amount",
		"host_refund_count",
		"host_refund_amount",
		"pos_sale_count",
		"pos_sale_amount",
		"pos_refund_count",
		"pos_refund_amount",
		"signature",
		"created_at",
		"process_settle",
	).Create(&entity)

	return result.Error
}

func (r Repo) UpdateFirstSettleDate(ctx context.Context, tx *gorm.DB, settltementID string, trxDate time.Time) error {
	result := tx.WithContext(ctx).Model(&Settlement{}).
		Where("settlement_id = ?", settltementID).Updates(&Settlement{FirstTrxTime: trxDate})

	return result.Error
}
