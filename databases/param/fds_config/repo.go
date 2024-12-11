package fdsconfig

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

func (r Repo) GetFdsAddress(ctx context.Context) (string, error) {
	var fdsConfig FdsConfig
	result := r.Db.WithContext(ctx).Select("fds_address").Find(&fdsConfig)

	return fdsConfig.FdsAddress, result.Error
}