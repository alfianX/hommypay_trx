package terminals

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

func (r Repo) CheckTidMid(ctx context.Context, tid string, mid string) (int64, error) {
	var count int64
	result := r.Db.WithContext(ctx).Model(&Terminals{}).Where("id = ? AND merchant_id = ?",
				tid, mid).Count(&count)
	
	return count, result.Error
}

func (r Repo) GetEmailMerchant(ctx context.Context, tid string, mid string) (string, error) {
	var terminal Terminals
	result := r.Db.WithContext(ctx).Select("email").Where("id = ? AND merchant_id = ?",
				tid, mid).Find(&terminal)

	return terminal.Email, result.Error
}

func (r Repo) GetBatch(ctx context.Context, tid string, mid string) (string, error) {
	var terminal Terminals
	result := r.Db.WithContext(ctx).Select("batch").Where("id = ? AND merchant_id = ?",
				tid, mid).Find(&terminal)

	return terminal.Batch, result.Error
}

func (r Repo) UpdateBatch(ctx context.Context, tx *gorm.DB, tid, mid, batch string) error {
	result := tx.WithContext(ctx).Model(&Terminals{}).Where(`
				id = ? AND merchant_id = ?`, tid, mid).Updates(&Terminals{
					Batch: batch,
				})

	return result.Error
}