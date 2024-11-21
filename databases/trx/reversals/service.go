package reversals

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

type SaveDataReversalParams struct {
	TransactionID   string
	TransactionType string
	Procode         string
	Mid             string
	Tid             string
	Amount          int64
	TransactionDate string
	Stan            string
	Trace           string
	Batch           string
	IsoRequest      string
}

type UpdateDataReversalParams struct {
	ID				int64
	ResponseCode	string
	IsoResponse		string
	RepeatCount		int64
}

type CheckDataReversalParams struct {
	Procode			string
	TID             string
	MID             string
	Amount          int64
	TransactionDate time.Time
	STAN            string
	Trace           string
	Batch			string
}

func (s Service) SaveDataReversal(ctx context.Context, params SaveDataReversalParams) error {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	trxDate, err := time.ParseInLocation("2006-01-02 15:04:05", params.TransactionDate, loc)
	if err != nil {
		return err
	}

	entity := Reversals{
		TransactionID: params.TransactionID,
		TransactionType: params.TransactionType,
		Procode: params.Procode,
		Mid: params.Mid,
		Tid: params.Tid,
		Amount: params.Amount,
		TransactionDate: trxDate,
		Stan: params.Stan,
		Trace: params.Trace,
		Batch: params.Batch,
		IsoRequest: params.IsoRequest,
		Flag: 70,
		CreatedAt: time.Now(),
	}

	err = s.repo.SaveDataReversal(ctx, &entity)
	
	return err
}

func (s Service) UpdateDataReversal(ctx context.Context, params UpdateDataReversalParams) error {
	var flag int64
	if params.ResponseCode == "00" {
		flag = 85
	}else{
		flag = 80
	}

	entity := Reversals{
		ID: params.ID,
		ResponseCode: params.ResponseCode,
		IsoResponse: params.IsoResponse,
		RepeatCount: params.RepeatCount,
		Flag: flag,
	}

	err := s.repo.UpdateDataReversal(ctx, &entity)

	return err
}

func (s *Service) CheckDataReversal(ctx context.Context, params CheckDataReversalParams) (int64, int64, string, error) {
	entity := Reversals{
		Procode: params.Procode,
		Mid: params.MID,
		Tid: params.TID,
		Amount: params.Amount,
		TransactionDate: params.TransactionDate,
		Stan: params.STAN,
		Trace: params.Trace,
		Batch: params.Batch,
	}
	id, flag, rcOrg, err := s.repo.CheckDataReversal(ctx, &entity)
	if err != nil {
		return 0, 0, "", err
	}

	return id, flag, rcOrg, nil
}
