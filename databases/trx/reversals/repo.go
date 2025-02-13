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
		"stan_issuer",
		"trace",
		"batch",
		"iso_request",
		"issuer_id",
		"flag",
		"created_at",
	).Create(&entity)

	return result.Error
}

func (r Repo) SaveDataReversalSettle(ctx context.Context, entity *Reversals) error {
	result := r.Db.WithContext(ctx).Select(
		"transaction_id",
		"transaction_type",
		"procode",
		"mid",
		"tid",
		"amount",
		"transaction_date",
		"stan",
		"stan_issuer",
		"trace",
		"batch",
		"iso_request",
		"issuer_id",
		"flag",
		"response_code_origin",
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

	return entity.ID, entity.Flag, entity.ResponseCodeOrigin, nil
}

func (r Repo) GetDataAutoReversal(ctx context.Context) ([]Reversals, error) {
	var allData []Reversals

	result := r.Db.WithContext(ctx).Select(
					"id",
					"transaction_id",
					"response_code_origin",
					"iso_request",
					"issuer_id",	
					"repeat_count",
				).Where(`transaction_type != ? AND response_code_origin IS NOT NULL AND flag = ?`, "41", 70).Find(&allData)

	return allData, result.Error
}

func (r Repo) UpdateFlagReversal(ctx context.Context, entity *Reversals) error {
	result := r.Db.WithContext(ctx).Model(&entity).Updates(&entity)

	return result.Error
}

func (r Repo) CreateAutoReversalLog(ctx context.Context, id int64) error {
	result := r.Db.WithContext(ctx).Exec(`INSERT INTO reversal_logs SELECT *, NOW() FROM reversals 
				WHERE id = ?`, id)

	return result.Error
}

func (r Repo) DeleteReversal(ctx context.Context, id int64) error {
	result := r.Db.WithContext(ctx).Exec(`DELETE FROM reversals WHERE id = ?`, id)

	return result.Error
}

func (r Repo) UpdateBackFlagReversal(ctx context.Context, entity *Reversals) error {
	result := r.Db.WithContext(ctx).Model(&entity).Updates(&entity)

	return result.Error
}

func (r Repo) GetDataSafReversal(ctx context.Context) ([]Reversals, error) {
	var allData []Reversals

	result := r.Db.WithContext(ctx).Select(
					"id",
					"transaction_id",
					"response_code_origin",
					"iso_request",
					"issuer_id",
					"repeat_count",
				).Where(`transaction_type = ? AND flag = ?`, "41", 70).Find(&allData)

	return allData, result.Error
}

