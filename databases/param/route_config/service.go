package routeconfig

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetUrlByEndPoint(ctx context.Context, endpoint string) (string, error) {
	url, err := s.repo.GetUrlByEndPoint(ctx, endpoint)
	if err != nil {
		return "", err
	}

	return url, nil
}