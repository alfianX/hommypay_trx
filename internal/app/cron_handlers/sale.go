package cronhandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	transactiondata "github.com/alfianX/hommypay_trx/databases/trx/transaction_data"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
)

func (cs *CronService) SaleTrx(params transactiondata.TransactionData) {
	
	type response struct {
		ResponseCode    string `json:"responseCode"`
		ApprovalCode  	string `json:"approvalCode"`
		// Signature	  string `json:"signature"`
		ISO8583 		string `json:"ISO8583"`
	}
	
	var req types.SaleRequest
	err := json.Unmarshal([]byte(params.DataRequest), &req)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Unmarshal data request : " + err.Error())
		fmt.Println(err)
	}

	ip, port, err := cs.hsmConfigService.GetHSMIpPort(context.Background())
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Get IP HSM : " + err.Error())
		fmt.Println(err)
	}

	zek, err := cs.keyConfigService.GetZEK(context.Background())
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Get ZEK : " + err.Error())
		fmt.Println(err)
	}

	lenPan, err := strconv.ParseInt(req.CardInformation.PAN[:4], 16, 0)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Parse int len pan : " + err.Error())
		fmt.Println(err)
	}
	pan, err := h.HSMDecrypt(ip+":"+port, zek, req.CardInformation.PAN[4:])
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Decrypt pan : " + err.Error())
		fmt.Println(err)
	}

	pan = pan[:lenPan]

	cardType, err := cs.binRangeService.GetCardTypeByPAN(context.Background(), pan)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Get card type : " + err.Error())
		fmt.Println(err)
	}

	pan = h.MaskPan(pan)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	trxDate, err := time.ParseInLocation("2006-01-0215:04:05", req.PaymentInformation.TransactionDate, loc)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Parse trx date : " + err.Error())
		fmt.Println(err)
	}

	trxParams := transactions.CreateTrxParams{
		TransactionID: params.TransactionID,
		Procode: req.PaymentInformation.Procode,
		Mid: req.PaymentInformation.MID,
		Tid: req.PaymentInformation.TID,
		CardType: cardType,
		Pan: pan,
		PanEnc: req.CardInformation.PAN,
		TrackData: req.CardInformation.TrackData2,
		EMVTag: req.CardInformation.EMVTag,
		Amount: req.PaymentInformation.Amount,
		TransactionDate: trxDate,
		Stan: req.PaymentInformation.STAN,
		Trace: req.PaymentInformation.Trace,
		Batch: req.PaymentInformation.Batch,
		TransMode: req.PosTerminal.TransMode,
		IsoRequest: req.ISO8583,
		IssuerID: params.IssuerID,
		Longitude: params.Longitude,
		Latitude: params.Latitude,
	}
	id, err := cs.transactionService.CreateSaleTrx(context.Background(), trxParams)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, 0, 0)
		h.ErrorLog("Cron - Save trx sale : " + err.Error())
		fmt.Println(err)
	}

	if params.DataResponse != "" {
		var responseCode string
		var iso8583 string
		var approvalCode string
		if params.DataResponse != "E6" {
			var res response
			err = json.Unmarshal([]byte(params.DataResponse), &res)
			if err != nil {
				cs.ErrorClear(context.Background(), params.ID, id, 1)
				h.ErrorLog("Cron - Unmarshal response sale : " + err.Error())
				fmt.Println(err)
			}

			responseCode = res.ResponseCode
			iso8583 = res.ISO8583
			approvalCode = res.ApprovalCode
		}else{
			responseCode = params.DataResponse
		}

		err = cs.transactionService.UpdateSaleTrx(context.Background(), transactions.UpdateSaleParams{
			ID: id,
			ResponseCode: responseCode,
			ISO8583Response: iso8583,
			ApprovalCode: approvalCode,
		})
		if err != nil {
			cs.ErrorClear(context.Background(), params.ID, id, 1)
			h.ErrorLog("Cron - Update trx sale : " + err.Error())
			fmt.Println(err)
		}
	}

	err = cs.transactionDataService.DeleteTrxData(context.Background(), params.ID)
	if err != nil {
		cs.ErrorClear(context.Background(), params.ID, id, 1)
		h.ErrorLog("Cron - Delete trx data : " + err.Error())
		fmt.Println(err)
	}
}

func (cs *CronService) ErrorClear(ctx context.Context, idTrxData int64, idTrx int64, afterSave int64) {
	cs.transactionDataService.UpdateFlagTrxDataBack(context.Background(), idTrxData)
	if afterSave > 0 {
		cs.transactionService.DeleteTrx(context.Background(), idTrx)
	}
}