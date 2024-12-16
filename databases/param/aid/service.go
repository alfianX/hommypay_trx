package aid

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetAppName(ctx context.Context, aid string) (string, error) {
	appName, err := s.repo.GetAppName(ctx, aid)
	if err != nil {
		return "", err
	}

	return appName, nil
}