package settlementdetails

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

type CreateSettleDetailParams struct {
	SettlementID        	int64
	TransactionID       	string
	TransactionType     	string
	Procode					string
	MID                 	string
	TID                 	string
	CardType            	string
	PAN                 	string
	PANEnc					string
	TrackData           	string
	EMVTag              	string
	Amount              	int64
	TransactionDate     	time.Time
	STAN                	string
	STANIssuer				string
	Rrn						string
	Trace               	string
	Batch					string
	TransMode				string
	BankCode				string
	ISO8583Request      	string
	ISO8583RequestIssuer	string
	ResponseCode        	string
	ResponseAt          	time.Time
	ApprovalCode        	string
	RefID               	string
	DE32					string
	DE33					string
	DE123					string
	ISO8583Response     	string
	ISO8583ResponseIssuer 	string
	IssuerID            	int64
	Signature				string
	VoidID              	string
	BatchUFlag          	int64
	CutOff					string
}

func (s Service) CreateSettleDetail(ctx context.Context, tx *gorm.DB, settleType string, params CreateSettleDetailParams) error {
	var status int64
	if settleType != "NORMAL" {
		if params.BatchUFlag == 2 {
			status = 1
		} else {
			status = 2
		}
	} else {
		status = 1
	}

	entity := SettlementDetails{
		SettlementID:        params.SettlementID,
		TransactionID:       params.TransactionID,
		TransactionType:     params.TransactionType,
		Procode: 			 params.Procode,
		Mid:                 params.MID,
		Tid:                 params.TID,
		CardType:            params.CardType,
		Pan:                 params.PAN,
		PanEnc:				 params.PANEnc,
		TrackData:           params.TrackData,
		EmvTag:              params.EMVTag,
		Amount:              params.Amount,
		TransactionDate:     params.TransactionDate,
		Stan:                params.STAN,
		StanIssuer: 		 params.STANIssuer,
		Trace:               params.Trace,
		Batch: 				 params.Batch,
		TransMode: 			 params.TransMode,
		BankCode: 			 params.BankCode,	
		IsoRequest:      	 params.ISO8583Request,
		IsoRequestIssuer:  	 params.ISO8583RequestIssuer,
		ResponseCode:        params.ResponseCode,
		ResponseAt:          params.ResponseAt,
		ApprovalCode:        params.ApprovalCode,
		RefID:               params.RefID,
		DE32: 				 params.DE32,
		DE33: 				 params.DE33,
		DE123: 				 params.DE123,
		IsoResponse:     	 params.ISO8583Response,
		IsoResponseIssuer: 	 params.ISO8583ResponseIssuer,
		IssuerID:            params.IssuerID,
		Signature: 			 params.Signature,
		Status:              status,
		VoidID:              params.VoidID,
		CutOff: 			 params.CutOff,
		CreatedAt: 			 time.Now(),		
	}

	err := s.repo.CreateSettleDetail(ctx, tx, &entity)
	
	return err
}