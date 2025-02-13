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

func (r Repo) CreateSettle(ctx context.Context, tx *gorm.DB, entity *Settlement) (int64, error) {
	result := tx.WithContext(ctx).Select(
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

	return entity.ID, result.Error
}

func (r Repo) UpdateFirstSettleDate(ctx context.Context, tx *gorm.DB, id int64, trxDate time.Time) error {
	result := tx.WithContext(ctx).Model(&Settlement{ID: id}).Updates(&Settlement{FirstTrxTime: trxDate})

	return result.Error
}