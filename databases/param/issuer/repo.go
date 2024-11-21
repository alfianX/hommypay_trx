package issuer

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

func (r Repo) GetUrlByIssuerID(ctx context.Context, issuerID int64) (string, int64, string, string, error) {
	var data Issuer
	data.ID = issuerID

	result := r.Db.WithContext(ctx).
				Select("issuer_name", "issuer_conn_type", "issuer_type", "issuer_service").First(&data)
	
	return data.IssuerName, data.IssuerConnType, data.IssuerType, data.IssuerService, result.Error
}