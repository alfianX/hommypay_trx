package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	uniqueid "github.com/albinj12/unique-id"
	"github.com/alfianX/hommypay_trx/databases/trx/settlement"
	settlementdetails "github.com/alfianX/hommypay_trx/databases/trx/settlement_details"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s service) Settlement(c *gin.Context) {
	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type response struct {
		Status             string `json:"status"`
		ResponseCode       string `json:"responseCode"`
		Message			   string `json:"message"`
		RefNo  	   		   string `json:"refNo"`
		Signature	  	   string `json:"signature"`
	}

	req := types.SettlementRequest{}

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

	// re := regexp.MustCompile(`\r?\n`)
	// dataRequest := re.ReplaceAllString(string(dataRequestByte), "")
	// dataRequest = strings.ReplaceAll(dataRequest, " ", "")

	// currentTime := time.Now()
	// gmtFormat := "15:04:05"
	// timeString := currentTime.Format(gmtFormat)
	// logMessage := fmt.Sprintf("[%s] - path:%s, method: %s,\n requestBody: %v", timeString, c.Request.URL.EscapedPath(), c.Request.Method, dataRequest)
	// h.HistoryLog(logMessage, "settlement")
	h.HistoryReqLog(c, dataRequestByte, dateString, timeString, "settlement")

	totalTransactionDB, totalAmountDB, err := s.transactionService.GetSettleTotal(c, mid, tid, settementType)
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

	saleCountDB, saleAmountDB, err := s.transactionService.GetSaleTotal(c, mid, tid, settementType)
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

	voidCountDB, voidAmountDB, err := s.transactionService.GetVoidTotal(c, mid, tid, settementType)
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

	id, err := s.settlementService.CreateSettle(c, settementType, settlement.CreateSettleParams{
		MID: mid,
		TID: tid,
		STAN: stan,
		Trace: trace,
		Batch: batch,
		RefNo: refNo,
		SettleDate: settleDate,
		TotalTransaction: totalTransaction,
		TotalAmount: totalAmount,
		SaleCount: saleCountDB,
		SaleAmount: saleAmountDB,
		VoidCount: voidCountDB,
		VoidAmount: voidAmountDB,
		PosSaleCount: saleCount,
		PosSaleAmount: saleAmount,
		PosVoidCount: voidCount,
		PosVoidAmount: voidAmount,
		Signature: signatureFinal,
	})
	if err != nil {
		h.ErrorLog("Save settle: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	dataTrx, err := s.transactionService.GetDataTrx(c, mid, tid)
	if err != nil {
		h.ErrorLog("Get data transactions: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	i := 1
	for _, data := range dataTrx {
		if i == 1 {
			err = s.settlementService.UpdateFirstSettleDate(c, id, data.TransactionDate)
			if err != nil {
				h.ErrorLog("Update first settle date: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
				return
			}
		}

		err = s.settementDetailService.CreateSettleDetail(c, settementType, settlementdetails.CreateSettleDetailParams{
			SettlementID: id,
			TransactionID: data.TransactionID,
			TransactionType: data.TransactionType,
			Procode: data.Procode,
			MID: data.Mid,
			TID: data.Tid,
			CardType: data.CardType,
			PAN: data.Pan,
			PANEnc: data.PanEnc,
			TrackData: data.TrackData,
			EMVTag: data.EmvTag,
			Amount: data.Amount,
			TransactionDate: data.TransactionDate,
			STAN: data.Stan,
			STANIssuer: data.StanIssuer,
			Rrn: data.Rrn,
			Trace: data.Trace,
			Batch: data.Batch,
			ISO8583Request: data.IsoRequest,
			ISO8583RequestIssuer: data.IsoRequestIssuer,
			ResponseCode: data.ResponseCode,
			ResponseAt: data.ResponseAt,
			ApprovalCode: data.ApprovalCode,
			RefID: data.ReffID,
			DE32: data.DE32,
			DE33: data.DE33,
			DE123: data.DE123,
			ISO8583Response: data.IsoResponse,
			ISO8583ResponseIssuer: data.IsoResponseIssuer,
			IssuerID: data.IssuerID,
			Signature: data.Signature,
			VoidID: data.VoidID,
			BatchUFlag: data.BatchUFlag,
			CutOff: data.SettledDate,
		})
		if err != nil {
			h.ErrorLog("Save settle detail: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		i++
	}

	err = s.transactionService.UpdateSettleFlag(c, mid, tid)
	if err != nil {
		h.ErrorLog("Update settle flag: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	err = s.transactionLogService.CreateLogTrx(c, mid, tid)
	if err != nil {
		h.ErrorLog("Save log trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	responseOK := response{
		Status: "SUCCESS",
		ResponseCode: "00",
		Message: "Approved",
		RefNo: refNo,
		Signature: signatureFinal,
	}

	dataResponseByte, err := json.Marshal(responseOK)
	if err != nil {
		h.ErrorLog("Marshal response : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	// dataResponse := re.ReplaceAllString(string(dataResponseByte), "")
	// dataResponse = strings.ReplaceAll(dataResponse, " ", "")

	// logMessage = fmt.Sprintf("\n respondStatus: %d, respondBody: %s\n", http.StatusOK, dataResponse)
	// h.HistoryLog(logMessage, "settlement")
	h.HistoryRespLog(dataResponseByte, dateString, timeString, "settlement")

	h.Respond(c, responseOK, http.StatusOK)
}