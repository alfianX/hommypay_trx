package fdsconfig

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetFdsAddress(ctx context.Context) (string, error) {
	fdsAddress, err := s.repo.GetFdsAddress(ctx)
	if err != nil {
		return "", err
	}

	return fdsAddress, nil
}