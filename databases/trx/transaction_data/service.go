package transactiondata

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

type TrxDataReqParams struct {
	TransactionID	string
	TransactionType	string
	DataReq			string
	IssuerID		int64
	Longitude		string
	Latitude		string
}

type TrxDataResParams struct {
	ID				int64
	DataRes			string
}

func (s Service) SaveTrxDataReq(ctx context.Context, params TrxDataReqParams) (int64, error) {
	entity := TransactionData{
		TransactionID: params.TransactionID,
		TransactionType: params.TransactionType,
		DataRequest: params.DataReq,
		IssuerID: params.IssuerID,
		Longitude: params.Longitude,
		Latitude: params.Latitude,
		CreatedAt: time.Now(),
		Flag: 70,
	}

	id, err := s.repo.SaveTrxDataReq(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) UpdateTrxDataRes(ctx context.Context, params TrxDataResParams) error {
	entity := TransactionData{
		ID: params.ID,
		DataResponse: params.DataRes,
		Flag: 80,
	}

	err := s.repo.UpdateTrxDataRes(ctx, &entity)

	return err
}

func (s Service) GetTrxData(ctx context.Context) ([]TransactionData, error) {
	data, err := s.repo.GetTrxData(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range data {
		entity := TransactionData{
			ID: row.ID,
			Flag: 85,
		}

		err = s.repo.UpdateFlagTrxData(ctx, &entity)
		if err != nil {
			return nil, err
		}	
	}

	return data, nil
}

func (s *Service) DeleteTrxData(ctx context.Context, id int64) error {
	entity := TransactionData{
		ID: id,
	}
	err := s.repo.DeleteTrxData(ctx, &entity)

	return err
}

func (s Service) UpdateFlagTrxDataBack(ctx context.Context, id int64) error {
	entity := TransactionData{
		ID: id,
		Flag: 80,
	}

	err := s.repo.UpdateFlagTrxData(ctx, &entity)

	return err
}

func (s Service) UpdateFlagTrxDataErr(ctx context.Context, id int64, responseCode string) error {
	entity := TransactionData{
		ID: id,
		DataResponse: responseCode,
		Flag: 80,
	}

	err := s.repo.UpdateFlagTrxData(ctx, &entity)

	return err
}