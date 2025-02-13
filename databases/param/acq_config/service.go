package acqconfig

import (
	"context"
	"strconv"
)

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetMaxAmount(ctx context.Context) (int64, error) {
	value, err := s.repo.GetMaxAmount(ctx)
	if err != nil {
		return 0, err
	}

	maxAmount, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return maxAmount, nil
}