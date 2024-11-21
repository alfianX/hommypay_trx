package settlement

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

type CreateSettleParams struct {
	MID              string
	TID              string
	STAN             string
	Trace            string
	Batch            string
	RefNo			 string
	SettleDate       string
	TotalTransaction int64
	TotalAmount      int64
	SaleCount        int64
	SaleAmount       int64
	VoidCount        int64
	VoidAmount       int64
	PosSaleCount     int64
	PosSaleAmount    int64
	PosVoidCount     int64
	PosVoidAmount    int64
}

func (s Service) CreateSettle(ctx context.Context, settleType string, params CreateSettleParams) (int64, error) {
	var subBatchNo string
	if settleType == "NORMAL" {
		subBatchNo = "00"
	} else {
		subBatchNo = "01"
	}

	entity := Settlement{
		Mid:              params.MID,
		Tid:              params.TID,
		Stan:             params.STAN,
		Trace:            params.Trace,
		Batch:            params.Batch,
		RefNo: 			  params.RefNo,
		SubBatchNo:       subBatchNo,
		SettleDate:       params.SettleDate,
		TotalTransaction: params.TotalTransaction,
		TotalAmount:      params.TotalAmount,
		HostSaleCount:    params.SaleCount,
		HostSaleAmount:   params.SaleAmount,
		PosSaleCount:     params.PosSaleCount,
		PosSaleAmount:    params.PosSaleAmount,
		CreatedAt: 		  time.Now(),
	}

	id, err := s.repo.CreateSettle(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s Service) UpdateFirstSettleDate(ctx context.Context, id int64, trxDate time.Time) error {
	err := s.repo.UpdateFirstSettleDate(ctx, id, trxDate)
	if err != nil {
		return err
	}

	return err
}