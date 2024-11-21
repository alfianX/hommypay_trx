package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	transactiondata "github.com/alfianX/hommypay_trx/databases/trx/transaction_data"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s *service) Reversal(c *gin.Context) {

	type send struct {
		TransactionID      string `json:"transactionID" binding:"required"`
		IssuerID           int64  `json:"issuerID" binding:"required"`
		PaymentInformation struct {
			Procode         string `json:"procode" binding:"required"`
			TID             string `json:"tid" binding:"required"`
			MID             string `json:"mid" binding:"required"`
			Amount          int64  `json:"amount" binding:"required"`
			Tip             int64  `json:"tip"`
			STAN            string `json:"stan" binding:"required"`
			Trace           string `json:"trace" binding:"required"`
			Batch           string `json:"batch"`
			TransactionDate string `json:"transactionDate" binding:"required"`
			KSN             string `json:"ksn"`
		} `json:"paymentInformation" binding:"required"`
		CardInformation struct {
			PAN        string `json:"pan" binding:"required"`
			Expiry     string `json:"expiry" `
			CardType   string `json:"cardType" binding:"required"`
			TrackData2 string `json:"trackData" binding:"required"`
			EMVTag     string `json:"emvTag"`
			PinBlock   string `json:"pinBlock"`
		} `json:"cardInformation" binding:"required"`
		PosTerminal struct {
			TransMode string `json:"transMode"`
			Code      string `json:"code" binding:"required"`
			KeyMode   int    `json:"keyMode" binding:"required"`
		} `json:"posTerminal"`
	}

	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type responseISO struct {
		Status        string `json:"status"`
		ResponseCode  string `json:"responseCode"`
		ISO8583       string `json:"ISO8583"`
	}

	type response struct {
		Status        string `json:"status"`
		ResponseCode  string `json:"responseCode"`
	}

	req := types.ReversalRequest{}

	err := h.Decode(c, &req)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]h.ErrorMsg, len(ve))
			for i, fe := range ve {
                out[i] = h.ErrorMsg{Field: fe.Field(), Message: h.GetErrorMsg(fe)}
            }
			
			// h.ErrorLog(err.Error())
			h.Respond(c, gin.H{"status": "INVALID_REQUEST", "ResponseCode": "I6", "Message": out}, http.StatusBadRequest)
			return
		}
		// h.ErrorLog(err.Error())
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I6", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	procode := req.PaymentInformation.Procode
	tid := req.PaymentInformation.TID
	mid := req.PaymentInformation.MID
	amount := req.PaymentInformation.Amount
	tip := req.PaymentInformation.Tip
	stan := req.PaymentInformation.STAN
	trace := req.PaymentInformation.Trace
	bacth := req.PaymentInformation.Batch
	transactionDate := req.PaymentInformation.TransactionDate
	var ksn string
	pan := req.CardInformation.PAN
	expiry := req.CardInformation.Expiry
	trackData2 := req.CardInformation.TrackData2
	emvTag := req.CardInformation.EMVTag
	var pinBlock string
	transMode := req.PosTerminal.TransMode
	code := req.PosTerminal.Code
	keyMode := req.PosTerminal.KeyMode
	ISO8583 := req.ISO8583
	var issuerID int64
	var lat string
	var long string
	var panEnc string
	var expiryEnc string
	var trackData2Enc string
	var emvTagEnc string
	var pinBlockEnc string
	var iso8583Enc string

	if c.GetHeader("X-LATITUDE") != "" {
		lat = c.GetHeader("X-LATITUDE")
	}

	if c.GetHeader("X-LONGITUDE") != "" {
		long = c.GetHeader("X-LONGITUDE")
	}

	if code == "02" {
		pinBlock = req.CardInformation.PinBlock
		ksn = req.PaymentInformation.KSN
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	trxDate, err := time.ParseInLocation("2006-01-02 15:04:05", transactionDate, loc)
	if err != nil {
		h.ErrorLog("Parse time: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
		return
	}

	trxID, issuerID, err := s.transactionService.CheckDataTrx(c, transactions.CheckDataTrxParams{
		Procode: procode,
		TID: tid,
		MID: mid,
		Amount: amount,
		TransactionDate: trxDate,
		STAN: stan,
		Trace: trace,
		Batch: bacth,
	})
	if err != nil {
		h.ErrorLog("Check data trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
		return
	}

	if trxID != "" {
		id, flag, rcOrg, err := s.reversalService.CheckDataReversal(c, reversals.CheckDataReversalParams{
			Procode: procode,
			TID: tid,
			MID: mid,
			Amount: amount,
			TransactionDate: trxDate,
			STAN: stan,
			Trace: trace,
			Batch: bacth,
		})
		if err != nil {
			h.ErrorLog("Check data reversal: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
			return
		}

		if id == 0 || (id != 0 && flag == 70 && rcOrg == "00") {
			var issuerName string
			var issuerConnType int64
			var issuerService string
			var cardType string
			issuerName, issuerConnType, cardType, issuerService, err = s.issuerService.GetUrlByIssuerID(c, issuerID)
			if err != nil {
				h.ErrorLog("Get url issuer: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			if issuerConnType == 0 {
				h.ErrorLog("Issuer conn type not found!")
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}
		
			if issuerService == "" {
				h.ErrorLog("Issuer service not found!")
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			ip, port, err := s.hsmConfigService.GetHSMIpPort(c)
			if err != nil {
				h.ErrorLog("Get ip address HSM: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			zek, err := s.keyConfigService.GetZEK(c)
			if err != nil {
				h.ErrorLog("Get ZEK: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			if pan != "" {
				panEnc, err = h.HSMEncrypt(ip+":"+port, zek, pan)
				if err != nil {
					h.ErrorLog("PAN encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.CardInformation.PAN = panEnc
			}

			if expiry != "" {
				expiryToEnc := strings.ReplaceAll(expiry, "/", "")
				expiryEnc, err = h.HSMEncrypt(ip+":"+port, zek, expiryToEnc)
				if err != nil {
					h.ErrorLog("Expiry encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.CardInformation.Expiry = expiryEnc
			}

			if trackData2 != "" {
				trackData2Enc, err = h.HSMEncrypt(ip+":"+port, zek, trackData2)
				if err != nil {
					h.ErrorLog("Trackdata2 encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.CardInformation.TrackData2 = trackData2Enc
			}

			if emvTag != "" {
				emvTagEnc, err = h.HSMEncrypt(ip+":"+port, zek, emvTag)
				if err != nil {
					h.ErrorLog("Emv tag encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.CardInformation.EMVTag = emvTagEnc
			}

			if pinBlock != "" {
				pinBlockEnc, err = h.HSMEncrypt(ip+":"+port, zek, pinBlock)
				if err != nil {
					h.ErrorLog("Pin block encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.CardInformation.PinBlock = pinBlockEnc
			}

			if ISO8583 != "" {
				iso8583Enc, err = h.HSMEncrypt(ip+":"+port, zek, ISO8583)
				if err != nil {
					h.ErrorLog("ISO req encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				req.ISO8583 = iso8583Enc
			}

			dataRequestByte, err := json.Marshal(req)
			if err != nil {
				h.ErrorLog("Marshal request : " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			re := regexp.MustCompile(`\r?\n`)
			dataRequest := re.ReplaceAllString(string(dataRequestByte), "")
			dataRequest = strings.ReplaceAll(dataRequest, " ", "")

			currentTime := time.Now()
			gmtFormat := "15:04:05"
			timeString := currentTime.Format(gmtFormat)
			logMessage := fmt.Sprintf("[%s] - path:%s, method: %s,\n requestBody: %v", timeString, c.Request.URL.EscapedPath(), c.Request.Method, dataRequest)
			h.HistoryLog(logMessage, "reversal")

			currentTime = time.Now().UTC()
			gmtFormat = "20060102150405"
			dateString := currentTime.Format(gmtFormat)
			transactionID := "TRX" + dateString + strconv.Itoa(time.Now().Nanosecond())[2:5]

			trxDataReqParams := transactiondata.TrxDataReqParams{
				TransactionID: transactionID,
				TransactionType: "41",
				DataReq: dataRequest,
				IssuerID: issuerID,
				Longitude: long,
				Latitude: lat,
			}

			id, err := s.transactionDataService.SaveTrxDataReq(c, trxDataReqParams)
			if err != nil {
				h.ErrorLog("Save trx: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			dataToSend := send{}
			dataToSend.TransactionID = transactionID
			dataToSend.IssuerID = issuerID
			dataToSend.PaymentInformation.Procode = procode
			dataToSend.PaymentInformation.TID = tid
			dataToSend.PaymentInformation.MID = mid
			dataToSend.PaymentInformation.Amount = amount
			dataToSend.PaymentInformation.Tip = tip
			dataToSend.PaymentInformation.TransactionDate = transactionDate
			dataToSend.PaymentInformation.STAN = stan
			dataToSend.PaymentInformation.KSN = ksn
			dataToSend.PaymentInformation.Trace = trace
			dataToSend.PaymentInformation.Batch = bacth
			dataToSend.CardInformation.PAN = pan
			dataToSend.CardInformation.Expiry = expiry
			dataToSend.CardInformation.CardType = cardType
			dataToSend.CardInformation.TrackData2 = trackData2
			dataToSend.CardInformation.EMVTag = emvTag
			dataToSend.CardInformation.PinBlock = pinBlock
			dataToSend.PosTerminal.TransMode = transMode
			dataToSend.PosTerminal.Code = code
			dataToSend.PosTerminal.KeyMode = keyMode

			payload, err := json.Marshal(dataToSend)
			if err != nil {
				s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
				h.ErrorLog("JSON marshal data send: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			var extResp map[string]interface{}
			var ISO8583Res string
			var dataResponse string
			var iso8583ResEnc string

			if issuerConnType == 1 {
				if ISO8583 != "" {
					logMessage := fmt.Sprintf("[%s] - Request: %s", timeString, iso8583Enc)
					h.IssuerLog(logMessage, issuerName)

					extResp, err = h.TcpSendToIssuer(c, s.config, ISO8583, issuerService)
					
				}else{
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "I6")
					// h.ErrorLog("ISO8583 empty!")
					h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I6", Message: "ISO8583 empty!"}, http.StatusBadRequest)
					return
				}
			}else if issuerConnType == 2 {
				logMessage := fmt.Sprintf("[%s] - Request: %s", timeString, dataRequest)
				h.IssuerLog(logMessage, issuerName)

				extResp, err = h.RestSendToIssuer(c, s.config, payload, issuerService)
			}

			if err != nil {
				if strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout"){
					err = s.transactionService.UpdateTOReversalFlag(c, trxID)
					if err != nil {
						s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
						h.ErrorLog("Update reversal TO flag: " + err.Error())
						h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
						return
					}
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "T1")
					// h.ErrorLog("request timeout")
					h.Respond(c, responseError{Status: "TIMEOUT", ResponseCode: "T1", Message: "request timeout"}, http.StatusConflict)
					return
				}else{
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "H1")
					h.ErrorLog("Send to microservice: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "H1", Message: "Host down"}, http.StatusConflict)
					return
				}
			}

			var responseCode string
			if extResp["responseCode"] != nil {
				responseCode = extResp["responseCode"].(string)
			}

			if responseCode == "00" {
				err = s.transactionService.UpdateReversalFlag(c, trxID)
				if err != nil {
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
					h.ErrorLog("Update reversal flag: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
			} else {
				err = s.transactionService.UpdateTOReversalFlag(c, trxID)
				if err != nil {
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
					h.ErrorLog("Update reversal TO flag: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
			}

			if extResp["ISO8583"] != nil {
				ISO8583Res = extResp["ISO8583"].(string)
				iso8583ResEnc, err = h.HSMEncrypt(ip+":"+port, zek, ISO8583Res)
				if err != nil {
					s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
					h.ErrorLog("ISO res encrypt: " + err.Error())
					h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
					return
				}
				extResp["ISO8583"] = iso8583ResEnc
			}
		
			dataResponseByte, err := json.Marshal(extResp)
			if err != nil {
				s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
				h.ErrorLog("Marshal response : " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}
		
			dataResponse = re.ReplaceAllString(string(dataResponseByte), "")
			dataResponse = strings.ReplaceAll(dataResponse, " ", "")
		
			if issuerConnType == 1 {
				logMessage = fmt.Sprintf("\n Response: %s\n", iso8583ResEnc)
				h.IssuerLog(logMessage, issuerName)
			}else if issuerConnType == 2 {
				logMessage = fmt.Sprintf("\n Response: %s\n", dataResponse)
				h.IssuerLog(logMessage, issuerName)
			}

			logMessage = fmt.Sprintf("\n respondStatus: %d, respondBody: %s\n", http.StatusOK, dataResponse)
			h.HistoryLog(logMessage, "reversal")

			err = s.transactionDataService.UpdateTrxDataRes(c, transactiondata.TrxDataResParams{
				ID: id,
				DataRes: dataResponse,
			})

			if err != nil {
				s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E6")
				h.ErrorLog("Update trx: " + err.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusConflict)
				return
			}

			responseStatus := "SUCCESS"
			if responseCode != "00" {
				responseStatus = "FAILURE"
			}

			if req.ISO8583 == "" {
				respons := response{}
				respons.Status = responseStatus
				respons.ResponseCode = responseCode
				h.Respond(c, respons, http.StatusOK)
			} else {
				respons := responseISO{}
				respons.Status = responseStatus
				respons.ResponseCode = responseCode
				respons.ISO8583 = ISO8583Res
				h.Respond(c, respons, http.StatusOK)
			}
		}else{
			respons := response{}
			respons.Status = "SUCCESS"
			respons.ResponseCode = "00"
			h.Respond(c, respons, http.StatusOK)
		}
	}else{
		respons := response{}
		respons.Status = "SUCCESS"
		respons.ResponseCode = "00"
		h.Respond(c, respons, http.StatusOK)
	}
}