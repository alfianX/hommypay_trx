package issuer

import "context"

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) GetUrlByIssuerID(ctx context.Context, issuerID int64) (string, int64, string, string, error) {
	issuerName, issuerConnType, issuerType, issuerService, err := s.repo.GetUrlByIssuerID(ctx, issuerID)
	if err != nil {
		return "", 0, "", "", err
	}

	return issuerName, issuerConnType, issuerType, issuerService, nil
}