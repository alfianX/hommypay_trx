package routeconfig

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

func (r Repo) GetUrlByEndPoint(ctx context.Context, endpoint string) (string, error) {
	var routeConfig RouteConfig
	result := r.Db.WithContext(ctx).Select("url").Where("endpoint = ? AND status = ?", endpoint, "1").Find(&routeConfig)

	return routeConfig.Url, result.Error
}