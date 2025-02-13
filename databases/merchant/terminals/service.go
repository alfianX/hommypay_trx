package terminals

import (
	"context"
	"fmt"
	"strconv"

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

func (s Service) CheckTidMid(ctx context.Context, tid, mid string) (int64, error) {
	count, err := s.repo.CheckTidMid(ctx, tid, mid)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s Service) GetEmailMerchant(ctx context.Context, tid string, mid string) (string, error) {
	email, err := s.repo.GetEmailMerchant(ctx, tid, mid)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (s Service) UpdateBatch(ctx context.Context, tx *gorm.DB, tid, mid string) error {
	var batch string
	
	batch, err := s.repo.GetBatch(ctx, tid, mid)
	if err != nil {
		return err
	}

	batchInt, err := strconv.Atoi(batch)
	if err != nil {
		return err
	}

	batchInt = batchInt + 1
	batchStr := fmt.Sprintf("%06d", batchInt)

	err = s.repo.UpdateBatch(ctx, tx, tid, mid, batchStr)

	return err
}