package reversals

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

func (r Repo) SaveDataReversal(ctx context.Context, entity *Reversals) error {
	result := r.Db.WithContext(ctx).Select(
		"transaction_id",
		"transaction_type",
		"procode",
		"mid",
		"tid",
		"amount",
		"transaction_date",
		"stan",
		"trace",
		"batch",
		"iso_request",
		"flag",
		"created_at",
	).Create(&entity)

	return result.Error
}

func (r Repo) UpdateDataReversal(ctx context.Context, entity *Reversals) error {
	result := r.Db.WithContext(ctx).Model(&entity).Updates(&entity)

	return result.Error
}

func (r Repo) CheckDataReversal(ctx context.Context, entity *Reversals) (int64, int64, string, error) {
	result := r.Db.WithContext(ctx).Select("id", "flag", "response_code_origin").Where(`procode = ? AND mid = ?
				AND tid = ? AND amount = ? AND transaction_date = ? AND stan = ? AND trace = ? AND batch = ?`,
				entity.Procode, entity.Mid, entity.Tid, entity.Amount, 
				entity.TransactionDate, entity.Stan, entity.Trace, entity.Batch,
				).First(&entity)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, 0, "", result.Error
	}

	return entity.ID, entity.Flag, entity.ResponseCodeOrg, nil
}