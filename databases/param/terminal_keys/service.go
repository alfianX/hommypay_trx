package terminalkeys

import (
	"context"
	"time"
)

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) SaveTPK(ctx context.Context, tid string, tpk string) error {
	count, err := s.repo.CheckKey(ctx, tid)
	if err != nil {
		return err
	}

	if count > 0 {
		err = s.repo.UpdateKey(ctx, tid, tpk)
	} else {
		entity := TerminalKeys{
			Tid: tid,
			Value: tpk,
			CreatedAt: time.Now(),
		}
		err = s.repo.CreateKey(ctx, &entity)
	}

	return err
}