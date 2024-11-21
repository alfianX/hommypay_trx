package transactionlogs

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

func (r Repo) CreateLogTrx(ctx context.Context, mid string, tid string) error {
	result := r.Db.WithContext(ctx).Exec(`INSERT INTO transaction_logs SELECT *, NOW() FROM transactions 
				WHERE mid = ? AND tid = ?`, mid, tid)

	return result.Error
}

func (r Repo) ClearTrx(ctx context.Context, mid string, tid string) error {
	result := r.Db.WithContext(ctx).Exec(`DELETE FROM transactions WHERE mid = ? AND tid = ?`, mid, tid)

	return result.Error
}