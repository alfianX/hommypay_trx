package hsmconfig

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetHSMIpPort(ctx context.Context) (string, string, error) {
	ip, port, err := s.repo.GetHSMIpPort(ctx)
	if err != nil {
		return "", "", err
	}

	return ip, port, nil
}