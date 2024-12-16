package suspectlist

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

type CreateSuspectParams struct {
	Mid    string
	Tid    string
	Trace  string
	Pan    string
	Date   string
	Status string
	Data   string
}

func (s Service) CreateSuspect(ctx context.Context, params CreateSuspectParams) error {
	entity := SuspectList{
		Mid: params.Mid,
		Tid: params.Tid,
		Trace: params.Trace,
		Pan: params.Pan,
		Date: params.Date,
		Status: params.Status,
		Data: params.Data,
	}

	err := s.repo.CreateSuspect(ctx, &entity)

	return err
}