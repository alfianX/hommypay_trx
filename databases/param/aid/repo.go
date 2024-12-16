package aid

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

func (r Repo) GetAppName(ctx context.Context, aid string) (string, error) {
	var aidList AidList
	result := r.Db.WithContext(ctx).Where("aid = ?", aid).First(&aidList)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return  "", result.Error
	}

	return aidList.AppName, nil
}