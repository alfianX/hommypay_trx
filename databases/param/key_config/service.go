package keyconfig

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetTMK(ctx context.Context) (string, error) {
	tmk, err := s.repo.GetTMK(ctx)
	if err != nil {
		return "", err
	}

	return tmk, nil
}

func (s Service) GetZEK(ctx context.Context) (string, error) {
	zek, err := s.repo.GetZEK(ctx)
	if err != nil {
		return "", err
	}

	return zek, nil
}