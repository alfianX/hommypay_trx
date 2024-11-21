package keyconfig

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

func (r Repo) GetTMK(ctx context.Context) (string, error) {
	var keyConfig KeyConfig
	result := r.Db.WithContext(ctx).Select("value").Where("key_type = ?", "TMK").Find(&keyConfig)

	return keyConfig.Value, result.Error
}

func (r Repo) GetZEK(ctx context.Context) (string, error) {
	var keyConfig KeyConfig
	result := r.Db.WithContext(ctx).Select("value").Where("key_type = ?", "ZEK").Find(&keyConfig)

	return keyConfig.Value, result.Error
}