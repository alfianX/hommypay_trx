package terminalkeys

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) CheckKey(ctx context.Context, tid string) (int64, error) {
	var count int64
	result := r.Db.WithContext(ctx).Model(&TerminalKeys{}).Where("tid = ? AND key_type = ?", tid, "TPK").Count(&count)

	return count, result.Error
}

func (r Repo) CreateKey(ctx context.Context, entity *TerminalKeys) error {
	result := r.Db.WithContext(ctx).Select("tid", "key_type", "value", "created_at").Create(&entity)

	return result.Error
}

func (r Repo) UpdateKey(ctx context.Context, tid string, tpk string) error {
	result := r.Db.WithContext(ctx).Model(&TerminalKeys{}).Where("tid = ? AND key_type = ?", tid, "TPK").Updates(TerminalKeys{Value: tpk, UpdatedAt: time.Now()})

	return result.Error
}