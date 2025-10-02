package routefds

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

func (r Repo) GetUrlFdsReject(ctx context.Context) ([]RouteFds, error) {
	var routeFds []RouteFds
	result := r.Db.WithContext(ctx).Select("url").Where("status = ? AND data = ?", 1, "Rejected").Find(&routeFds)

	return routeFds, result.Error
}

func (r Repo) GetUrlFdsSuspect(ctx context.Context) ([]RouteFds, error) {
	var routeFds []RouteFds
	result := r.Db.WithContext(ctx).Select("url").Where("status = ? AND data = ?", 1, "Suspect").Find(&routeFds)

	return routeFds, result.Error
}
