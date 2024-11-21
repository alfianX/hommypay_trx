package terminals

import "context"

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