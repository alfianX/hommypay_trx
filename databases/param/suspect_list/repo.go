package suspectlist

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

func (r Repo) CreateSuspect(ctx context.Context, entity *SuspectList) error {
	result := r.Db.WithContext(ctx).Select(
		"id",
		"MID",
		"TID",
		"trace",
		"PAN",
		"date",
		"status",
		"data",
	).Create(&entity)

	return result.Error
}