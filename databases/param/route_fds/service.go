package routefds

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetUrlFdsReject(ctx context.Context) ([]RouteFds, error) {
	routeFds, err := s.repo.GetUrlFdsReject(ctx)
	if err != nil {
		return []RouteFds{}, err
	}

	return routeFds, nil
}

func (s Service) GetUrlFdsSuspect(ctx context.Context) ([]RouteFds, error) {
	routeFds, err := s.repo.GetUrlFdsSuspect(ctx)
	if err != nil {
		return []RouteFds{}, err
	}

	return routeFds, nil
}
