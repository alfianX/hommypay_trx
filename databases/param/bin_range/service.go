package binrange

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetUrlByPAN(ctx context.Context, pan, cardType string) (int64, string, int64, string, error) {
	issuerID, issuerName, issuerConnType, issuerService, err := s.repo.GetUrlByPAN(ctx, pan, cardType)
	if err != nil {
		return 0, "", 0, "", err
	}

	return issuerID, issuerName, issuerConnType, issuerService, nil
}

func (s Service) GetCardTypeByPAN(ctx context.Context, pan string) (string, error) {
	cardType, err := s.repo.GetCardTypeByPAN(ctx, pan)
	if err != nil {
		return "", err
	}

	return cardType, nil
}