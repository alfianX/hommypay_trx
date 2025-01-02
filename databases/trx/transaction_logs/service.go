package transactionlogs

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) CreateLogTrx(ctx context.Context, mid, tid, batch string) error {
	tx := s.repo.Db.Begin()
	
	defer tx.Rollback()

	err := s.repo.CreateLogTrx(ctx, mid, tid, batch)
	if err != nil {
		return err
	}

	err = s.repo.ClearTrx(ctx, mid, tid, batch)
	if err != nil {
		return err
	}

	tx.Commit()

	return err
}