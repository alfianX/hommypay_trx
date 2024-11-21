package handlers

import (
	"encoding/json"
	"net/http"

	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func (s service) Actions(c *gin.Context) {
	type responseError struct {
		Status       string `json:"status"`
		ResponseCode string `json:"responseCode"`
		Message      string `json:"message"`
	}

	req, err := c.GetRawData()
	if err != nil {
		h.ErrorLog(err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusBadRequest)
		return
	}

	action := c.Param("action")
	url, err := s.routeConfig.GetUrlByEndPoint(c, action)
	if err != nil {
		h.ErrorLog("Get url microservice: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusInternalServerError)
		return
	}

	headerMap := make(map[string]string)
    for key, values := range c.Request.Header {
        if len(values) > 0 {
            headerMap[key] = values[0]
        }
    }

	client := resty.New()
	response, err := client.R().
		SetHeaders(headerMap).
		SetBody(req).
		Post(url)

	if err != nil {
		h.ErrorLog("Send to microservice: " + err.Error())
		h.Respond(c, responseError{Status: "SERVER_FAILED", ResponseCode: "E6", Message: "Service Malfunction"}, http.StatusInternalServerError)
		return
	}

	var res interface{}
	json.Unmarshal(response.Body(), &res)

	h.Respond(c, res, response.StatusCode())
}