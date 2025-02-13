package transactionlogs

import (
	"context"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) CreateLogTrx(ctx context.Context, tx *gorm.DB, mid, tid, batch string) error {

	err := s.repo.CreateLogTrx(ctx, tx, mid, tid, batch)
	if err != nil {
		return err
	}

	err = s.repo.ClearTrx(ctx, tx, mid, tid, batch)
	if err != nil {
		return err
	}

	return err
}