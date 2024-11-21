package transactions

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

type CreateTrxParams struct {
	TransactionID     string
	Procode			  string
	Mid               string
	Tid	              string
	CardType          string
	Pan               string
	PanEnc            string
	TrackData         string
	EMVTag            string
	Amount            int64
	TransactionDate   time.Time
	Stan              string
	Trace             string
	Batch			  string
	TransMode		  string
	IsoRequest	      string
	IssuerID          int64
	Longitude         string
	Latitude          string
}

type UpdateSaleParams struct {
	ID              int64
	ResponseCode    string
	ISO8583Response string
	ApprovalCode    string
}

type UpdateVoidParams struct {
	ID              int64
	TransactionID   string
	ResponseCode    string
	ISO8583Response string
	ApprovalCode    string
	SaleID          int64
}

type CheckDataTrxParams struct {
	Procode			string
	TID             string
	MID             string
	Amount          int64
	TransactionDate time.Time
	STAN            string
	Trace           string
	Batch			string
}

func (s Service) CreateSaleTrx(ctx context.Context, params CreateTrxParams) (int64, error) {
	entity := Transactions{
		TransactionID: params.TransactionID,
		TransactionType: "01",
		Procode: params.Procode,
		Mid: params.Mid,
		Tid: params.Tid,
		CardType: params.CardType,
		Pan: params.Pan,
		PanEnc: params.PanEnc,
		TrackData: params.TrackData,
		EmvTag: params.EMVTag,
		Amount: params.Amount,
		TransactionDate: params.TransactionDate,
		Stan: params.Stan,
		Trace: params.Trace,
		Batch: params.Batch,
		TransMode: params.TransMode,
		IsoRequest: params.IsoRequest,
		IssuerID: params.IssuerID,
		Status: 1,
		Longitude: params.Longitude,
		Latitude: params.Latitude,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateTrx(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) CreateReversalTrx(ctx context.Context, params CreateTrxParams) (int64, error) {
	entity := Transactions{
		TransactionID: params.TransactionID,
		TransactionType: "41",
		Procode: params.Procode,
		Mid: params.Mid,
		Tid: params.Tid,
		CardType: params.CardType,
		Pan: params.Pan,
		PanEnc: params.PanEnc,
		TrackData: params.TrackData,
		EmvTag: params.EMVTag,
		Amount: params.Amount,
		TransactionDate: params.TransactionDate,
		Stan: params.Stan,
		Trace: params.Trace,
		Batch: params.Batch,
		TransMode: params.TransMode,
		IsoRequest: params.IsoRequest,
		IssuerID: params.IssuerID,
		Status: 1,
		Longitude: params.Longitude,
		Latitude: params.Latitude,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateTrx(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) CreateVoidTrx(ctx context.Context, params CreateTrxParams) (int64, error) {
	entity := Transactions{
		TransactionID: params.TransactionID,
		TransactionType: "31",
		Procode: params.Procode,
		Mid: params.Mid,
		Tid: params.Tid,
		CardType: params.CardType,
		Pan: params.Pan,
		PanEnc: params.PanEnc,
		TrackData: params.TrackData,
		EmvTag: params.EMVTag,
		Amount: params.Amount,
		TransactionDate: params.TransactionDate,
		Stan: params.Stan,
		Trace: params.Trace,
		Batch: params.Batch,
		TransMode: params.TransMode,
		IsoRequest: params.IsoRequest,
		IssuerID: params.IssuerID,
		Status: 1,
		Longitude: params.Longitude,
		Latitude: params.Latitude,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateTrx(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) UpdateSaleTrx(ctx context.Context, params UpdateSaleParams) error {
	
	entity := Transactions{
		ID:              params.ID,
		ResponseCode:    params.ResponseCode,
		IsoResponse: 	 params.ISO8583Response,
		ApprovalCode:    params.ApprovalCode,
	}

	err := s.repo.UpdateTrx(ctx, &entity)
	if err != nil {
		return err
	}
	
	return err
}

func (s Service) UpdateReversalTrx(ctx context.Context, params UpdateSaleParams) error {
	
	entity := Transactions{
		ID:              params.ID,
		ResponseCode:    params.ResponseCode,
		IsoResponse: 	 params.ISO8583Response,
		ApprovalCode:    params.ApprovalCode,
	}

	err := s.repo.UpdateTrx(ctx, &entity)
	if err != nil {
		return err
	}
	
	return err
}

func (s Service) UpdateVoidTrx(ctx context.Context, params UpdateVoidParams) error {
	tx := s.repo.Db.Begin()

	defer tx.Rollback()

	entity := Transactions{
		ID:              params.ID,
		ResponseCode:    params.ResponseCode,
		IsoResponse: 	 params.ISO8583Response,
		ApprovalCode:    params.ApprovalCode,
	}

	err := s.repo.UpdateTrx(ctx, &entity)
	if err != nil {
		return err
	}

	if params.ResponseCode == "00" {
		err = s.repo.UpdateVoidID(ctx, params.TransactionID, params.SaleID)
		if err != nil {
			return err
		}
	}
	
	tx.Commit()
	return err
}

func (s Service) GetSettleTotal(ctx context.Context, mid string, tid string, settleType string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var err error

	if settleType == "NORMAL" {
		totalTransaction, totalAmount, err = s.repo.GetSettleTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	} else {
		totalTransaction, totalAmount, err = s.repo.GetSettleBatchTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	}
	return totalTransaction, totalAmount, nil
}

func (s Service) GetSaleTotal(ctx context.Context, mid string, tid string, settleType string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var err error

	if settleType == "NORMAL" {
		totalTransaction, totalAmount, err = s.repo.GetSaleTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	} else {
		totalTransaction, totalAmount, err = s.repo.GetSaleBatchTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	}
	return totalTransaction, totalAmount, nil
}

func (s Service) GetVoidTotal(ctx context.Context, mid string, tid string, settleType string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var err error

	if settleType == "NORMAL" {
		totalTransaction, totalAmount, err = s.repo.GetVoidTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	} else {
		totalTransaction, totalAmount, err = s.repo.GetVoidBatchTotal(ctx, mid, tid)
		if err != nil {
			return 0, 0, err
		}
	}
	return totalTransaction, totalAmount, nil
}

func (s Service) GetDataTrx(ctx context.Context, mid string, tid string) ([]Transactions, error) {
	data, err := s.repo.GetDataTrx(ctx, mid, tid)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s Service) UpdateReversal(ctx context.Context, id int64) error {

	entity := Transactions{
		ID: id,
	}

	err := s.repo.UpdateReversal(ctx, &entity)
	if err != nil {
		return err
	}

	return err
}

func (s Service) UpdateSettleFlag(ctx context.Context, mid string, tid string) error {

	err := s.repo.UpdateSettleFlag(ctx, mid, tid)
	if err != nil {
		return err
	}

	return err
}

func (s Service) UpdateReversalFlag(ctx context.Context, trxID string) error {

	tx := s.repo.Db.Begin()
	
	defer tx.Rollback()

	err := s.repo.UpdateReversalVoidID(ctx, trxID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateReversalFlag(ctx, trxID, 1)
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}

func (s Service) UpdateTOReversalFlag(ctx context.Context, trxID string) error {
	err := s.repo.UpdateReversalFlag(ctx, trxID, 2)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) CheckDataTrx(ctx context.Context, params CheckDataTrxParams) (string, int64, error) {
	entity := Transactions{
		Procode: 		 params.Procode,
		Tid:             params.TID,
		Mid:             params.MID,
		Amount:          params.Amount,
		TransactionDate: params.TransactionDate,
		Stan:            params.STAN,
		Trace:           params.Trace,
		Batch: 			 params.Batch,
	}

	trxID, issuerID, err := s.repo.CheckDataTrx(ctx, &entity)
	if err != nil {
		return "", 0, err
	}

	return trxID, issuerID, err
}

func (s Service) CheckDataTrxV2(ctx context.Context, params CheckDataTrxParams) (string, int64, error) {
	entity := Transactions{
		Procode: 		 params.Procode,
		Tid:             params.TID,
		Mid:             params.MID,
		Amount:          params.Amount,
		TransactionDate: params.TransactionDate,
		Stan:            params.STAN,
		Trace:           params.Trace,
	}

	trxID, issuerID, err := s.repo.CheckDataTrxV2(ctx, &entity)
	if err != nil {
		return "", 0, err
	}

	return trxID, issuerID, err
}

func (s Service) CheckBatchDataTrx(ctx context.Context, params CheckDataTrxParams) (int64, error) {
	
	entity := Transactions{
		TransactionType: "01",
		Procode: 		 params.Procode,
		Tid:             params.TID,
		Mid:             params.MID,
		Amount:          params.Amount,
		TransactionDate: params.TransactionDate,
		Stan:            params.STAN,
		Trace:           params.Trace,
	}

	id, err := s.repo.CheckBatchDataTrx(ctx, &entity)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s Service) UpdateBatchFlag(ctx context.Context, id int64) error {
	entity := Transactions{
		ID: id,
	}

	err := s.repo.UpdateBatchFlag(ctx, &entity)
	if err != nil {
		return err
	}

	return err
}

func (s Service) GetTraceNoByIdTrx(ctx context.Context, trxID string) (string, error) {
	trace, err := s.repo.GetTraceNoByIdTrx(ctx, trxID)
	if err != nil {
		return "", err
	}

	return trace, nil
}

func (s Service) GetDataByTrxID(ctx context.Context, trxID string) (Transactions, error) {
	data, err := s.repo.GetDataByTrxID(ctx, trxID)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (s *Service) DeleteTrx(ctx context.Context, id int64) error {
	entity := Transactions{
		ID: id,
	}
	err := s.repo.DeleteTrx(ctx, &entity)

	return err
}