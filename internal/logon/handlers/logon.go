package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/gin-gonic/gin"
)

func (s service) Logon(c *gin.Context) {
	type request struct {
		DeviceInformation struct {
			MID string `json:"mid" binding:"required"`
			TID string `json:"tid" binding:"required"`
		} `json:"deviceInformation" binding:"required"`
	}

	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	type response struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message		 string `json:"message"`
		Key          string `json:"key"`
	}

	req := request{}

	err := h.Decode(c, &req)
	if err != nil {
		h.ErrorLog(err.Error())
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I0", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	timeFormat := "15:04:05"
	timeString := currentTime.Format(timeFormat)
	dateFormat := "20060102"
	dateString := currentTime.Format(dateFormat)

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

	dataRequestByte, err := json.Marshal(req)
	if err != nil {
		h.ErrorLog("Marshal request : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusBadRequest)
		return
	}

	// re := regexp.MustCompile(`\r?\n`)
	// dataRequest := re.ReplaceAllString(string(dataRequestByte), "")
	// dataRequest = strings.ReplaceAll(dataRequest, " ", "")

	// currentTime := time.Now()
	// gmtFormat := "15:04:05"
	// dateString := currentTime.Format(gmtFormat)
	// logMessage := fmt.Sprintf("[%s] - path:%s, method: %s,\n requestBody: %v", dateString, c.Request.URL.EscapedPath(), c.Request.Method, dataRequest)
	// h.HistoryLog(logMessage, "logon")
	h.HistoryReqLog(c, dataRequestByte, dateString, timeString, "logon")

	count, err := s.terminalService.CheckTidMid(c, req.DeviceInformation.TID, req.DeviceInformation.MID)
	if err != nil {
		h.ErrorLog("Check TID MID: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	if count == 0 {
		h.Respond(c, responseError{Status: "INVALID_REQUEST", ResponseCode: "I7", Message: "TID MID not registered!"}, http.StatusBadRequest)
		return
	}

	tmk, err := s.keyConfigService.GetTMK(c)
	if err != nil {
		h.ErrorLog("Get TMK: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	command := "0000HC" + tmk + ";XU0"
	len := len(command)
	lenHex := strings.ToUpper(fmt.Sprintf("%04x", len))
	message := lenHex + hex.EncodeToString([]byte(command))

	responseHSM, err := h.SendMessageToHsm(ip+":"+port, message)
	if err != nil {
		h.ErrorLog("Send to HSM: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	resByte, err := hex.DecodeString(responseHSM)
	if err != nil {
		h.ErrorLog("Decode response: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	responseCode := string(resByte[8:10])
	if responseCode != "00" {
		h.ErrorLog("HSM response code " + responseCode)
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}
	twk := string(resByte[11:43])
	tpk := string(resByte[43:76])

	err = s.terminalKeysService.SaveTPK(c, req.DeviceInformation.TID, tpk)
	if err != nil {
		h.ErrorLog("Save TPK: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	keyEnc, err := h.HSMEncrypt(ip+":"+port, zek, hex.EncodeToString([]byte(twk)))
	if err != nil {
		h.ErrorLog("Key encrypt: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusConflict)
		return
	}

	resp := response{}
	resp.Status = "SUCCESS"
	resp.ResponseCode = "00"
	resp.Message = "Approved"
	resp.Key = keyEnc

	dataResponseByte, err := json.Marshal(resp)
	if err != nil {
		h.ErrorLog("Marshal response : " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E1", Message: "Service Acq Malfunction"}, http.StatusBadRequest)
		return
	}

	// dataResponse := re.ReplaceAllString(string(dataResponseByte), "")
	// dataResponse = strings.ReplaceAll(dataResponse, " ", "")

	// logMessage = fmt.Sprintf("\n respondStatus: %d, respondBody: %s\n", http.StatusOK, dataResponse)
	// h.HistoryLog(logMessage, "logon")
	h.HistoryRespLog(dataResponseByte, dateString, timeString, "logon")

	resp.Key = hex.EncodeToString([]byte(twk))

	h.Respond(c, resp, http.StatusOK)
}