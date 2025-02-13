package acqconfig

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

func (r Repo) GetMaxAmount(ctx context.Context) (string, error) {
	var acqConfig AcqConfig
	result := r.Db.WithContext(ctx).Where("name = ?", "MAX_AMOUNT").First(&acqConfig)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return  "", result.Error
	}

	return acqConfig.Value, nil
}