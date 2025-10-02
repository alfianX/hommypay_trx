package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s service) BatchUpload(c *gin.Context) {
	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type response struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	req := types.BatchUploadRequest{}

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

	procode := req.PaymentInformation.Procode
	tid := req.PaymentInformation.TID
	mid := req.PaymentInformation.MID
	amount := req.PaymentInformation.Amount
	transactionDate := req.PaymentInformation.TransactionDate
	stan := req.PaymentInformation.STAN
	trace := req.PaymentInformation.Trace
	batch := req.PaymentInformation.Batch

	currentTime := time.Now()
	timeFormat := "15:04:05"
	timeString := currentTime.Format(timeFormat)
	dateFormat := "20060102"
	dateString := currentTime.Format(dateFormat)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	trxDate, err := time.ParseInLocation("2006-01-02 15:04:05", transactionDate, loc)
	if err != nil {
		h.ErrorLog("Parse time: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	dataRequestByte, err := json.Marshal(req)
	if err != nil {
		h.ErrorLog("Marshal request : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	h.HistoryReqLog(c, dataRequestByte, dateString, timeString, "batch_upload")

	trxID, err := s.transactionService.CheckBatchDataTrx(c, transactions.CheckDataTrxParams{
		Procode:         procode,
		TID:             tid,
		MID:             mid,
		Amount:          amount,
		TransactionDate: trxDate,
		STAN:            stan,
		Trace:           trace,
		Batch:           batch,
	})
	if err != nil {
		h.ErrorLog("Check data trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if trxID == "" {
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I1", Message: "Trx not found"}, http.StatusConflict)
		return
	}

	err = s.transactionService.UpdateBatchFlag(c, trxID)
	if err != nil {
		h.ErrorLog("Update flag batch upload: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	responseOK := response{
		Status:       "SUCCESS",
		ResponseCode: "00",
		Message:      "Approved",
	}

	dataResponseByte, err := json.Marshal(responseOK)
	if err != nil {
		h.ErrorLog("Marshal response : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	h.HistoryRespLog(dataResponseByte, dateString, timeString, "batch_upload")

	h.Respond(c, responseOK, http.StatusOK)
}
