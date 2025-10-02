package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	uniqueid "github.com/albinj12/unique-id"
	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	"github.com/alfianX/hommypay_trx/databases/trx/settlement"
	settlementdetails "github.com/alfianX/hommypay_trx/databases/trx/settlement_details"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type contextKey string

const RouteKey contextKey = "current_route"

func (s service) Settlement(c *gin.Context) {
	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type response struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
		RefNo        string `json:"refNo"`
		Signature    string `json:"signature"`
	}

	c.Set("current_route", "settlement")

	req := types.SettlementRequest{}

	path := c.Request.URL.Path

	err := h.Decode(c, &req)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]h.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = h.ErrorMsg{Field: fe.Field(), Message: h.GetErrorMsg(fe)}
			}

			// h.ErrorLog(err.Error())
			h.Respond(c, gin.H{"status": "INVALID_REQUEST", "ResponseCode": "I0", "Message": out}, http.StatusBadRequest)
			return
		}
		// h.ErrorLog(err.Error())
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I0", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	var processSettle string
	settementType := req.SettlementType
	tid := req.PaymentInformation.TID
	mid := req.PaymentInformation.MID
	stan := req.PaymentInformation.STAN
	trace := req.PaymentInformation.Trace
	batch := req.PaymentInformation.Batch
	settleDate := req.PaymentInformation.SettleDate
	totalTransaction := req.OrderInformation.TotalTransaction
	totalAmount := req.OrderInformation.TotalAmount
	saleCount := req.OrderInformation.SaleCount
	saleAmount := req.OrderInformation.SaleAmount
	voidCount := req.OrderInformation.VoidCount
	voidAmount := req.OrderInformation.VoidAmount

	if path == "/settlement" {
		processSettle = "NORMAL"
	} else if path == "/manual-settlement" {
		processSettle = "MANUAL"
	}

	countSettle, err := s.transactionService.CheckDataSettle(c, transactions.CheckDataSettleParams{
		TID:   tid,
		MID:   mid,
		Batch: batch,
	})
	if err != nil {
		h.ErrorLog("Check data settle : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if countSettle == 0 {
		responseOK := response{
			Status:       "SUCCESS",
			ResponseCode: "00",
			Message:      "Approved",
			RefNo:        "",
			Signature:    "",
		}

		h.Respond(c, responseOK, http.StatusOK)
		return
	}

	currentTime := time.Now()
	timeFormat := "15:04:05"
	timeString := currentTime.Format(timeFormat)
	dateFormat := "20060102"
	dateString := currentTime.Format(dateFormat)

	dataRequestByte, err := json.Marshal(req)
	if err != nil {
		h.ErrorLog("Marshal request : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	h.HistoryReqLog(c, dataRequestByte, dateString, timeString, "settlement")

	totalTransactionDB, totalAmountDB, err := s.transactionService.GetSettleTotal(c, mid, tid, batch, settementType)
	if err != nil {
		h.ErrorLog("Get settle total: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if settementType != "END" {
		if totalTransactionDB != totalTransaction || totalAmountDB != totalAmount {
			h.Respond(c, responseError{Status: "INVALID_SETTLE", ResponseCode: "95", Message: "Settle not match"}, http.StatusBadRequest)
			return
		}
	}

	saleCountDB, saleAmountDB, err := s.transactionService.GetSaleTotal(c, mid, tid, batch, settementType)
	if err != nil {
		h.ErrorLog("Get sale total: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if settementType != "END" {
		if saleCountDB != saleCount || saleAmountDB != saleAmount {
			h.Respond(c, responseError{Status: "INVALID_SETTLE", ResponseCode: "95", Message: "Settle not match"}, http.StatusBadRequest)
			return
		}
	}

	voidCountDB, voidAmountDB, err := s.transactionService.GetVoidTotal(c, mid, tid, batch, settementType)
	if err != nil {
		h.ErrorLog("Get void total: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if settementType != "END" {
		if voidCountDB != voidCount || voidAmountDB != voidAmount {
			h.Respond(c, responseError{Status: "INVALID_SETTLE", ResponseCode: "95", Message: "Settle not match"}, http.StatusBadRequest)
			return
		}
	}

	refNo, err := uniqueid.Generateid("n", 12)
	if err != nil {
		h.ErrorLog("Create ref no: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	email, err := s.terminalService.GetEmailMerchant(c, req.PaymentInformation.TID, req.PaymentInformation.MID)
	if err != nil {
		h.ErrorLog("Get email merchant : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	signatureFinal, err := h.CreateSignature(req.PaymentInformation.TID, req.PaymentInformation.MID, email, req.PaymentInformation.SettleDate, req.PaymentInformation.Trace, refNo)
	if err != nil {
		h.ErrorLog("Create signature: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	tx := s.dbTrx.Begin()
	if tx.Error != nil {
		h.ErrorLog("Db transaction begin : " + tx.Error.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	dateTimeFormat := "20060102150405"
	dateTimeString := currentTime.Format(dateTimeFormat)
	settlementID := "STL" + dateTimeString + strconv.Itoa(time.Now().Nanosecond())[2:5]

	err = s.settlementService.CreateSettle(c, tx, settementType, settlement.CreateSettleParams{
		SettlementID:     settlementID,
		MID:              mid,
		TID:              tid,
		STAN:             stan,
		Trace:            trace,
		Batch:            batch,
		RefNo:            refNo,
		SettleDate:       settleDate,
		TotalTransaction: totalTransaction,
		TotalAmount:      totalAmount,
		SaleCount:        saleCountDB,
		SaleAmount:       saleAmountDB,
		VoidCount:        voidCountDB,
		VoidAmount:       voidAmountDB,
		PosSaleCount:     saleCount,
		PosSaleAmount:    saleAmount,
		PosVoidCount:     voidCount,
		PosVoidAmount:    voidAmount,
		Signature:        signatureFinal,
		ProcessSettle:    processSettle,
	})
	if err != nil {
		tx.Rollback()
		h.ErrorLog("Save settle: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	dataTrx, err := s.transactionService.GetDataTrx(c, mid, tid, batch)
	if err != nil {
		h.ErrorLog("Get data transactions: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	i := 1
	for _, data := range dataTrx {
		if i == 1 {
			err = s.settlementService.UpdateFirstSettleDate(c, tx, settlementID, data.TransactionDate)
			if err != nil {
				tx.Rollback()
				h.ErrorLog("Update first settle date: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
				return
			}
		}

		if settementType == "END" && data.BatchUFlag == 1 {
			err := s.reversalService.SaveDataReversalSettle(c, reversals.SaveDataReversalParams{
				TransactionID:   data.TransactionID,
				TransactionType: data.TransactionType,
				Procode:         data.Procode,
				Mid:             data.Mid,
				Tid:             data.Tid,
				Amount:          data.Amount,
				TransactionDate: data.TransactionDate,
				Stan:            data.Stan,
				StanIssuer:      data.StanIssuer,
				Trace:           data.Trace,
				Batch:           data.Batch,
				IsoRequest:      data.IsoRequest,
				IssuerID:        data.IssuerID,
				ResponseCodeOrg: data.ResponseCode,
			})

			if err != nil {
				tx.Rollback()
				h.ErrorLog("Save auto reversal : " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
				return
			}
		}

		err = s.settementDetailService.CreateSettleDetail(c, tx, settementType, settlementdetails.CreateSettleDetailParams{
			SettlementID:    settlementID,
			TransactionID:   data.TransactionID,
			TransactionType: data.TransactionType,
			Procode:         data.Procode,
			MID:             data.Mid,
			TID:             data.Tid,
			CardType:        data.CardType,
			PAN:             data.Pan,
			PANEnc:          data.PanEnc,
			EMVTag:          data.EmvTag,
			Amount:          data.Amount,
			TransactionDate: data.TransactionDate,
			STAN:            data.Stan,
			STANIssuer:      data.StanIssuer,
			Rrn:             data.Rrn,
			Trace:           data.Trace,
			Batch:           data.Batch,
			TransMode:       data.TransMode,
			BankCode:        data.BankCode,
			DE43:            data.DE43,
			ResponseCode:    data.ResponseCode,
			ResponseAt:      data.ResponseAt,
			ApprovalCode:    data.ApprovalCode,
			RefID:           data.ReffID,
			DE32:            data.DE32,
			DE33:            data.DE33,
			DE123:           data.DE123,
			IssuerID:        data.IssuerID,
			Signature:       data.Signature,
			VoidID:          data.VoidID,
			BatchUFlag:      data.BatchUFlag,
			CutOff:          data.SettledDate,
		})
		if err != nil {
			tx.Rollback()
			h.ErrorLog("Save settle detail: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		i++
	}

	err = s.transactionService.UpdateSettleFlag(c, tx, mid, tid, batch)
	if err != nil {
		tx.Rollback()
		h.ErrorLog("Update settle flag: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	err = s.transactionLogService.CreateLogTrx(c, tx, mid, tid, batch)
	if err != nil {
		tx.Rollback()
		h.ErrorLog("Save log trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	txMerchant := s.dbMerchant.Begin()
	if txMerchant.Error != nil {
		h.ErrorLog("Db merchant transaction begin : " + tx.Error.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	err = s.terminalService.UpdateBatch(c, txMerchant, tid, mid)
	if err != nil {
		txMerchant.Rollback()
		tx.Rollback()
		h.ErrorLog("Update Batch: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		txMerchant.Rollback()
		h.ErrorLog("Commit : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if err := txMerchant.Commit().Error; err != nil {
		tx.Rollback()
		txMerchant.Rollback()
		h.ErrorLog("Commit db merchant : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	responseOK := response{
		Status:       "SUCCESS",
		ResponseCode: "00",
		Message:      "Approved",
		RefNo:        refNo,
		Signature:    signatureFinal,
	}

	dataResponseByte, err := json.Marshal(responseOK)
	if err != nil {
		h.ErrorLog("Marshal response : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	h.HistoryRespLog(dataResponseByte, dateString, timeString, "settlement")

	h.Respond(c, responseOK, http.StatusOK)
}
