package transactions

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) CheckData(ctx context.Context, entity *Transactions) (int64, error) {
	var count int64

	result := r.Db.WithContext(ctx).Model(&entity).
	Where(`transaction_type = ? AND procode = ? AND mid = ? AND tid = ? AND amount = ? AND 
			transaction_date = ? AND stan = ? AND trace = ? AND batch = ? AND response_code IS NOT NULL`, 
			entity.TransactionType, entity.Procode,	entity.Mid, entity.Tid, entity.Amount, 
			entity.TransactionDate, entity.Stan, entity.Trace, entity.Batch).Count(&count)
	
	return count, result.Error
}

func (r Repo) CheckStan(ctx context.Context, entity *Transactions, dateNow string) (int64, error) {
	var count int64

	result := r.Db.WithContext(ctx).Model(&entity).
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND stan = ? AND 
					DATE_FORMAT(transaction_date, '%Y-%m-%d') = ? AND response_code IS NOT NULL`, entity.TransactionType, entity.Mid,
						entity.Tid, entity.Stan, dateNow).Count(&count)

	return count, result.Error
}

func (r Repo) CreateTrx(ctx context.Context, entity *Transactions) (int64, error) {
	result := r.Db.WithContext(ctx).Select(
		"transaction_id",
		"transaction_type",
		"procode",
		"mid",
		"tid",
		"card_type",
		"pan",
		"pan_enc",
		"track_data",
		"emv_tag",
		"amount",
		"transaction_date",
		"stan",
		"trace",
		"batch",
		"trans_mode",
		"bank_code",
		"iso_request",
		"issuer_id",
		"status",
		"longitude",
		"latitude",
		"created_at",
	).Create(&entity)

	return entity.ID, result.Error
}

func (r Repo) UpdateTrx(ctx context.Context, entity *Transactions) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{ID: entity.ID}).Updates(&Transactions{
		ResponseCode: entity.ResponseCode,
		ResponseAt: time.Now(),
		IsoResponse: entity.IsoResponse,
		ApprovalCode: entity.ApprovalCode,
		Signature: entity.Signature,
		Status: 2,
		UpdatedAt: time.Now(),
	})

	return result.Error
}

func (r Repo) UpdateVoidID(ctx context.Context, voidId string, idTrx int64) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{ID: idTrx}).Updates(&Transactions{VoidID: voidId})

	return result.Error
}

func (r Repo) UpdateReversal(ctx context.Context, entity *Transactions) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{ID: entity.ID}).Updates(&Transactions{ReversalFlag: 1})

	return result.Error
}

func (r Repo) GetSettleTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND void_id IS NULL AND 
						batch_u_flag = ?`,
						"01", mid, tid, batch, 2, "00", 1, 0, 1).Find(&trx)
	
	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetSettleBatchTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND void_id IS NULL AND 
						batch_u_flag = ?`, "01", mid, tid, batch, 2, "00", 1, 0, 2).Find(&trx)

	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetSaleTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND void_id IS NULL AND 
						batch_u_flag = ?`, "01", mid, tid, batch, 2, "00", 1, 0, 1).Find(&trx)

	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetVoidTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND 
						batch_u_flag = ?`, "31", mid, tid, batch, 2, "00", 1, 0, 1).Find(&trx)

	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetSaleBatchTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND 
						batch_u_flag = ?`, "01", mid, tid, batch, 2, "00", 1, 0, 2).Find(&trx)

	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetVoidBatchTotal(ctx context.Context, mid, tid, batch string) (int64, int64, error) {
	var totalTransaction int64
	var totalAmount int64
	var trx []Transactions

	result := r.Db.WithContext(ctx).Select("amount").
				Where(`transaction_type = ? AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
						response_code = ? AND reversal_flag != ? AND settle_flag = ? AND 
						batch_u_flag = ?`, "31", mid, tid, batch, 2, "00", 1, 0, 2).Find(&trx)

	for _, rows := range trx {
		totalAmount = totalAmount + rows.Amount
		totalTransaction++
	}

	return totalTransaction, totalAmount, result.Error
}

func (r Repo) GetDataTrx(ctx context.Context, mid, tid, batch string) ([]Transactions, error) {
	var allData []Transactions

	result := r.Db.WithContext(ctx).Order("transaction_date asc").
				Where(`transaction_type IN ('01', '31') AND mid = ? AND tid = ? AND batch = ? AND status = ? AND 
				response_code = ? AND reversal_flag != ? AND settle_flag = ?`,
				mid,tid,batch,2,"00",1,0).Find(&allData)
	
	return allData, result.Error
}

func (r Repo) UpdateSettleFlag(ctx context.Context, tx *gorm.DB, mid, tid, batch string) error {
	result := tx.WithContext(ctx).Model(&Transactions{}).
				Where(`mid = ? AND tid = ? AND batch = ? AND settle_flag = ?`, mid, tid, batch, 0).
				Updates(&Transactions{SettleFlag: 1, SettledAt: time.Now()})
	
	return result.Error
}

func (r Repo) CheckDataTrx(ctx context.Context, entity *Transactions) (string, int64, string, error) {
	result := r.Db.WithContext(ctx).Select("transaction_id", "issuer_id", "bank_code").
				Where(`procode = ? AND mid = ? AND tid = ? AND amount = ?
				AND transaction_date = ? AND stan = ? AND trace = ? AND batch = ? AND status = ? 
				AND response_code = ? AND reversal_flag != ?`,
				entity.Procode, entity.Mid, entity.Tid, entity.Amount, 
				entity.TransactionDate, entity.Stan, entity.Trace, entity.Batch, 2, "00", 1).
				Find(&entity)
	
	return entity.TransactionID, entity.IssuerID, entity.BankCode, result.Error
}

func (r Repo) CheckDataTrxV2(ctx context.Context, entity *Transactions) (string, int64, error) {
	result := r.Db.WithContext(ctx).Select("transaction_id", "issuer_id").
				Where(`procode = ? AND mid = ? AND tid = ? AND amount = ?
				AND trace = ? AND status = ? AND response_code = ?`,
				entity.Procode, entity.Mid, entity.Tid, entity.Amount, entity.Trace, 2, "00").
				Find(&entity)
	
	return entity.TransactionID, entity.IssuerID, result.Error
}

func (r Repo) UpdateReversalFlag(ctx context.Context, trxId string, flag int64) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{}).
				Where(`transaction_id = ? AND status = ? AND response_code = ?`,
				trxId, 2, "00").Updates(&Transactions{ReversalFlag: flag, UpdatedAt: time.Now()})

	return result.Error
}

func (r Repo) UpdateReversalFlagTO(ctx context.Context, trxId string, flag int64) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{}).
				Where(`transaction_id = ?`,trxId).
				Updates(&Transactions{ReversalFlag: flag, UpdatedAt: time.Now()})

	return result.Error
}

func (r Repo) UpdateReversalVoidID(ctx context.Context, trxId string) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{}).
				Where(`void_id = ? AND status = ? AND response_code = ? AND reversal_flag != ?`,
				trxId, 2, "00", 1).Update("void_id", nil)

	return result.Error
}

func (r Repo) CheckBatchDataTrx(ctx context.Context, entity *Transactions) (int64, error) {
	result := r.Db.WithContext(ctx).Select("id").
				Where(`transaction_type = ? AND procode = ? AND mid = ? AND tid = ? AND amount = ? AND 
				transaction_date = ? AND trace = ? AND batch = ? AND
				status = ? AND response_code = ? AND reversal_flag != ?`,
				entity.TransactionType, entity.Procode, entity.Mid, entity.Tid, entity.Amount, 
				entity.TransactionDate, entity.Trace, entity.Batch, 2, "00", 1).
				Find(&entity)

	return entity.ID, result.Error
}

func (r Repo) UpdateBatchFlag(ctx context.Context, entity *Transactions) error {
	result := r.Db.WithContext(ctx).Model(&Transactions{ID: entity.ID}).
				Updates(&Transactions{BatchUFlag: 2, UpdatedAt: time.Now()})

	return result.Error
}

func (r Repo) GetTraceNoByIdTrx(ctx context.Context, trxId string) (string, error) {
	var data Transactions

	result := r.Db.WithContext(ctx).Select("trace").
				Where(`transaction_id = ? AND status = ? AND response_code = ? AND reversal_flag != ? AND 
				settle_flag = ? AND void_id IS NULL`, trxId, 2, "00", 1, 0).Find(&data)
				
	return data.Trace, result.Error
}

func (r Repo) GetDataByTrxID(ctx context.Context, trxID string) (Transactions, error) {
	var data Transactions

	result := r.Db.WithContext(ctx).
				Where(`transaction_id = ? AND status = ? AND response_code = ? AND 
				reversal_flag != ? AND settle_flag = ? AND void_id IS NULL`, trxID, 2, "00", 1, 0).Find(&data)
	
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return data, result.Error
	}

	return data, nil
}



func (r Repo) DeleteTrx(ctx context.Context, entity *Transactions) error {
	result := r.Db.WithContext(ctx).Model(&entity).Delete(&entity)

	return result.Error
}

func (r Repo) CheckDataSettle(ctx context.Context, entity *Transactions) (int64, error) {
	var count int64

	result := r.Db.WithContext(ctx).Model(&entity).
				Where(`mid = ? AND tid = ? AND batch = ?`, 
						entity.Mid, entity.Tid, entity.Batch).Count(&count)

	return count, result.Error
}

func (r Repo) GetTrxByTrxID(ctx context.Context, trxID string) (Transactions, error) {
	var data Transactions

	result := r.Db.WithContext(ctx).
				Where(`transaction_id = ?`, trxID).Find(&data)
	
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return data, result.Error
	}

	return data, nil
}