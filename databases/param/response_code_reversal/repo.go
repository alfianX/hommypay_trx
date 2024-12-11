package responsecodereversal

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

func (r Repo) GetRC(ctx context.Context, rc string) (string, error) {
	var responseCodeReversal ResponseCodeReversal
	result := r.Db.WithContext(ctx).Where("code = ?", rc).First(&responseCodeReversal)

	return responseCodeReversal.Description, result.Error
}