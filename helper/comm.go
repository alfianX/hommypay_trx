package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/pkg/iso"
	"github.com/gin-gonic/gin"
)

func RestSendToIssuer(c *gin.Context, cnf configs.Config, payload []byte, issuerService string) (map[string]interface{}, error) {
	var resp *http.Response

	exReq, err := http.NewRequest("POST", issuerService+"/sale", bytes.NewReader(payload))
	if err != nil {
		return nil, errors.New("Prepare send to microservice: " + err.Error())
	}

	exReq.Header = c.Request.Header

	client := &http.Client{
		Timeout: time.Duration(cnf.TimeoutTrx) * time.Second,
	}
	resp, err = client.Do(exReq)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var extResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&extResp)
	if err != nil {
		return nil, errors.New("Decode response: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Response: " + extResp["message"].(string))
	}

	return extResp, nil
}

func TcpSendToIssuer(c *gin.Context, cnf configs.Config, payload string, issuerService string) (map[string]interface{}, error) {
	var approvalCode string
	var responseCode string
	var isoResponse string
	var extResp map[string]interface{}

	ret, msg := SendHost(issuerService, payload, cnf.TimeoutTrx)
	if ret != 0 {
		return nil, errors.New("Send to host: " + msg)
	}

	lenIso := len([]rune(msg))

	if lenIso > 8 {
		if msg[4:6] == "60" {
			DE := iso.Parse(msg[14:])

			if DE[38] != "" {
				approvalCode = DE[38]
			}

			if DE[39] != "" {
				responseCode = DE[39]
			}

			isoResponse = msg
		}
	}
	
	extResp = map[string]interface{}{
		"responseCode": responseCode,
		"approvalCode": approvalCode,
		"ISO8583": strings.ToUpper(isoResponse),
	}

	return extResp, nil
}