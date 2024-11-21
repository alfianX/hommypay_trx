package transactiondata

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

func (r Repo) SaveTrxDataReq(ctx context.Context, entity *TransactionData) (int64, error) {
	result := r.Db.WithContext(ctx).Select(
		"transaction_id",
		"transaction_type",
		"data_request",
		"issuer_id",
		"longitude",
		"latitude",
		"created_at",
		"flag",
	).Create(&entity)

	return entity.ID, result.Error
}

func (r Repo) UpdateTrxDataRes(ctx context.Context, entity *TransactionData) error {
	result := r.Db.WithContext(ctx).Model(&entity).Updates(&entity)

	return result.Error
}

func (r Repo) GetTrxData(ctx context.Context) ([]TransactionData, error) {
	var allData []TransactionData

	result := r.Db.WithContext(ctx).Order("created_at asc").
				Where(`flag = ?`, 80).Find(&allData)

	return allData, result.Error
}

func (r Repo) UpdateFlagTrxData(ctx context.Context, entity *TransactionData) error {
	result := r.Db.WithContext(ctx).Model(&entity).Updates(&entity)

	return result.Error
}

func (r Repo) DeleteTrxData(ctx context.Context, entity *TransactionData) error {
	result := r.Db.WithContext(ctx).Model(&entity).Delete(&entity)

	return result.Error
}