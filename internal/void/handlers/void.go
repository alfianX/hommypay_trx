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
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (s service) Void(c *gin.Context) {
	
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
			PAN        string `json:"pan"`
			Expiry     string `json:"expiry" `
			CardType   string `json:"cardType"`
			TrackData2 string `json:"trackData"`
			EMVTag     string `json:"emvTag"`
			PinBlock   string `json:"pinBlock"`
		} `json:"cardInformation"`
		PosTerminal struct {
			TransMode 	string `json:"transMode"`
			Code    	string `json:"code"`
			KeyMode 	int    `json:"keyMode"`
		} `json:"posTerminal"`
	}

	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type responseISO struct {
		Status             string `json:"status"`
		ResponseCode       string `json:"responseCode"`
		Message      string `json:"message"`
		TransactionID 	   string `json:"transactionID"`
		ApprovalCode       string `json:"approvalCode"`
		Signature	  string `json:"signature"`
		ISO8583 string `json:"ISO8583"`
	}

	type response struct {
		Status             string `json:"status"`
		ResponseCode       string `json:"responseCode"`
		Message      string `json:"message"`
		TransactionID 	   string `json:"transactionID"`
		ApprovalCode  	   string `json:"approvalCode"`
		Signature	  string `json:"signature"`
	}

	req := types.VoidRequest{}

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

	transactionIdSale := req.TransactionID
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
	var issuerName string
	var issuerConnType int64
	var issuerService string
	var cardType string
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

	data, err := s.transactionService.GetDataByTrxID(c, transactionIdSale)
	if err != nil {
		h.ErrorLog("Get data sale: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if data.ID == 0 {
		// h.ErrorLog("Transaction not found!")
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I1", Message: "Transaction not found!"}, http.StatusBadRequest)
		return
	}

	if data.IssuerID == 0 {
		h.ErrorLog("Issuer not found!")
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Issuer not found!"}, http.StatusConflict)
		return
	}

	if data.IssuerID == 99 {
		if ISO8583 == "" {
			// h.ErrorLog("ISO8583 empty!")
			h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I2", Message: "ISO8583 empty!"}, http.StatusBadRequest)
			return
		}
	}

	issuerID = data.IssuerID

	issuerName, issuerConnType, cardType, issuerService, err = s.issuerService.GetUrlByIssuerID(c, issuerID)
	if err != nil {
		h.ErrorLog("Get url issuer service: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if issuerConnType == 0 {
		h.ErrorLog("Issuer conn type not found!")
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if issuerService == "" {
		h.ErrorLog("Issuer service not found!")
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if data.Trace != trace {
		// h.ErrorLog("Trace not found!")
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I4", Message: "Trace not found!"}, http.StatusBadRequest)
		return
	}

	ip, port, err := s.hsmConfigService.GetHSMIpPort(c)
	if err != nil {
		h.ErrorLog("Get ip address HSM: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	zek, err := s.keyConfigService.GetZEK(c)
	if err != nil {
		h.ErrorLog("Get ZEK: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if pan != "" {
		panEnc, err = h.HSMEncrypt(ip+":"+port, zek, pan)
		if err != nil {
			h.ErrorLog("PAN encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.CardInformation.PAN = panEnc
	}

	if expiry != "" {
		expiryToEnc := strings.ReplaceAll(expiry, "/", "")
		expiryEnc, err = h.HSMEncrypt(ip+":"+port, zek, expiryToEnc)
		if err != nil {
			h.ErrorLog("Expiry encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.CardInformation.Expiry = expiryEnc
	}

	if trackData2 != "" {
		trackData2Enc, err = h.HSMEncrypt(ip+":"+port, zek, trackData2)
		if err != nil {
			h.ErrorLog("Trackdata2 encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.CardInformation.TrackData2 = trackData2Enc
	}

	if emvTag != "" {
		emvTagEnc, err = h.HSMEncrypt(ip+":"+port, zek, emvTag)
		if err != nil {
			h.ErrorLog("Emv tag encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.CardInformation.EMVTag = emvTagEnc
	}

	if pinBlock != "" {
		pinBlockEnc, err = h.HSMEncrypt(ip+":"+port, zek, pinBlock)
		if err != nil {
			h.ErrorLog("Pin block encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.CardInformation.PinBlock = pinBlockEnc
	}

	if ISO8583 != "" {
		iso8583Enc, err = h.HSMEncrypt(ip+":"+port, zek, ISO8583)
		if err != nil {
			h.ErrorLog("ISO req encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		req.ISO8583 = iso8583Enc
	}

	dataRequestByte, err := json.Marshal(req)
	if err != nil {
		h.ErrorLog("Marshal request : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	re := regexp.MustCompile(`\r?\n`)
	dataRequest := re.ReplaceAllString(string(dataRequestByte), "")
	dataRequest = strings.ReplaceAll(dataRequest, " ", "")

	currentTime := time.Now()
	gmtFormat := "15:04:05"
	timeString := currentTime.Format(gmtFormat)
	logMessage := fmt.Sprintf("[%s] - path:%s, method: %s,\n requestBody: %v", timeString, c.Request.URL.EscapedPath(), c.Request.Method, dataRequest)
	h.HistoryLog(logMessage, "sale")

	currentTime = time.Now().UTC()
	gmtFormat = "20060102150405"
	dateString := currentTime.Format(gmtFormat)
	transactionID := "TRX" + dateString + strconv.Itoa(time.Now().Nanosecond())[2:5]

	trxDataReqParams := transactiondata.TrxDataReqParams{
		TransactionID: transactionID,
		TransactionType: "31",
		DataReq: dataRequest,
		IssuerID: issuerID,
		Longitude: long,
		Latitude: lat,
	}

	id, err := s.transactionDataService.SaveTrxDataReq(c, trxDataReqParams)
	if err != nil {
		h.ErrorLog("Save trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
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
		s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
		h.ErrorLog("JSON marshal data send: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
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
			s.transactionDataService.UpdateFlagTrxDataErr(c, id, "I2")
			h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I2", Message: "ISO8583 empty!"}, http.StatusBadRequest)
			return
		}
	}else if issuerConnType == 2 {
		logMessage := fmt.Sprintf("[%s] - Request: %s", timeString, dataRequest)
		h.IssuerLog(logMessage, issuerName)

		extResp, err = h.RestSendToIssuer(c, s.config, payload, issuerService)
	}
 
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout"){
			errRvrsl := s.AutoReversal(c, req, transactionID, issuerID, "")
			if errRvrsl != nil {
				s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
				h.ErrorLog("Save data reversal: " + errRvrsl.Error())
				h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
				return
			}
			s.transactionDataService.UpdateFlagTrxDataErr(c, id, "T0")
			// h.ErrorLog("request timeout")
			h.Respond(c, responseError{Status: "TIMEOUT", ResponseCode: "T0", Message: "request timeout"}, http.StatusConflict)
			return
		}else{
			s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E0")
			h.ErrorLog("Send to microservice: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E0", Message: "Link down!"}, http.StatusConflict)
			return
		}
	}

	var responseCode string
	if extResp["responseCode"] != nil {
		responseCode = extResp["responseCode"].(string)
	}
	var message string
	if extResp["message"] != nil {
		message = extResp["message"].(string)
	}
	var approvalCode string
	if extResp["approvalCode"] != nil {
		approvalCode = extResp["approvalCode"].(string)
	}

	if extResp["ISO8583"] != nil {
		ISO8583Res = extResp["ISO8583"].(string)
		iso8583ResEnc, err = h.HSMEncrypt(ip+":"+port, zek, ISO8583Res)
		if err != nil {
			s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
			h.ErrorLog("ISO res encrypt: " + err.Error())
			h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
			return
		}
		extResp["ISO8583"] = iso8583ResEnc
	}

	dataResponseByte, err := json.Marshal(extResp)
	if err != nil {
		s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
		h.ErrorLog("Marshal response : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
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

	email, err := s.terminalService.GetEmailMerchant(c, req.PaymentInformation.TID, req.PaymentInformation.MID)
	if err != nil {
		s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
		h.ErrorLog("Get email merchant : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}
	
	signatureFinal, err := h.CreateSignature(req.PaymentInformation.TID, req.PaymentInformation.MID, email, req.PaymentInformation.TransactionDate, req.PaymentInformation.Trace, approvalCode)
	if err != nil {
		s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
		h.ErrorLog("Create signature: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}


	logMessage = fmt.Sprintf("\n respondStatus: %d, respondBody: %s\n", http.StatusOK, dataResponse)
	h.HistoryLog(logMessage, "sale")

	err = s.transactionDataService.UpdateTrxDataRes(c, transactiondata.TrxDataResParams{
		ID: id,
		DataRes: dataResponse,
	})

	if err != nil {
		s.transactionDataService.UpdateFlagTrxDataErr(c, id, "E1")
		h.ErrorLog("Update trx: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	// select { 
	// 	case <-c.Request.Context().Done(): 
	// 	fmt.Println("Client disconnected") 
	// 	return 
	// default: 
	// 	fmt.Fprintln(c.Writer, "Client is still connected") 
	// }
	testTimeout := s.config.TestTimeout
	timeInt := s.config.Timeout
	if testTimeout == 1 {
		time.Sleep(time.Duration(timeInt) * time.Second)
	}

	responseStatus := "SUCCESS"
	if responseCode != "00" {
		responseStatus = "FAILURE"
	}

	if req.ISO8583 == "" {
		respons := response{}
		respons.Status = responseStatus
		respons.ResponseCode = responseCode
		respons.Message = message
		respons.TransactionID = transactionID
		respons.ApprovalCode = approvalCode
		respons.Signature = signatureFinal
		h.Respond(c, respons, http.StatusOK)
	} else {
		respons := responseISO{}
		respons.Status = responseStatus
		respons.ResponseCode = responseCode
		respons.Message = message
		respons.TransactionID = transactionID
		respons.ApprovalCode = approvalCode
		respons.Signature = signatureFinal
		respons.ISO8583 = ISO8583Res
		h.Respond(c, respons, http.StatusOK)
	}
}

func (s service) AutoReversal(c *gin.Context, req types.VoidRequest, trxId string, issuerID int64, rcOrg string) error {
	err := s.reversalService.SaveDataReversal(c, reversals.SaveDataReversalParams{
		TransactionID: trxId,
		TransactionType: "31",
		Procode: req.PaymentInformation.Procode,
		Mid: req.PaymentInformation.MID,
		Tid: req.PaymentInformation.TID,
		Amount: req.PaymentInformation.Amount,
		TransactionDate: req.PaymentInformation.TransactionDate,
		Stan: req.PaymentInformation.STAN,
		Trace: req.PaymentInformation.Trace,
		Batch: req.PaymentInformation.Batch,
		IsoRequest: req.ISO8583,
		IssuerID: issuerID,
		ResponseCodeOrg: rcOrg,
	})

	return err
}