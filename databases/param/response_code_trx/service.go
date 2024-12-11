package responsecodetrx

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetRC(ctx context.Context, rc string) (string, error) {
	description, err := s.repo.GetRC(ctx, rc)
	if err != nil {
		return "", err
	}

	return description, nil
}