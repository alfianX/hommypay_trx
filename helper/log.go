package helper

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorLog(logMessage string) {
	currentTime := time.Now()
	gmtFormat := "20060102"
	dateString := currentTime.Format(gmtFormat)
	filename := fmt.Sprintf("../log/error_log/err_%s.log", dateString)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	currentTime = time.Now()
	gmtFormat = "15:04:05"
	dateString = currentTime.Format(gmtFormat)
	logRequest := fmt.Sprintf("[%s] - %s\n\n", dateString, logMessage)
	file.WriteString(logRequest)
}

func HistoryReqLog(c *gin.Context, dataRequestByte []byte, dateString, timeString, name string) {
	if gin.Mode() == "debug" {
		re := regexp.MustCompile(`\r?\n`)
		dataRequest := re.ReplaceAllString(string(dataRequestByte), "")
		dataRequest = strings.ReplaceAll(dataRequest, " ", "")
		dataRequest = strings.ReplaceAll(dataRequest, "cardInformation", "ci")
		dataRequest = strings.ReplaceAll(dataRequest, "aid", "a")
		dataRequest = strings.ReplaceAll(dataRequest, "pan", "p")
		dataRequest = strings.ReplaceAll(dataRequest, "expiry", "ex")
		dataRequest = strings.ReplaceAll(dataRequest, "trackData", "td")
		dataRequest = strings.ReplaceAll(dataRequest, "emvTag", "et")
		dataRequest = strings.ReplaceAll(dataRequest, "pinBlock", "pb")
		dataRequest = strings.ReplaceAll(dataRequest, "ISO8583", "i")

		logMessage := fmt.Sprintf("[%s] - path:%s, method: %s,\n requestBody: %v", timeString, c.Request.URL.EscapedPath(), c.Request.Method, dataRequest)

		filename := fmt.Sprintf("../log/history_log/%s_%s.log", name, dateString)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		file.WriteString(logMessage)
	}
}

func HistoryRespLog(dataResponseByte []byte, dateString, timeString, name string) {
	if gin.Mode() == "debug" {
		re := regexp.MustCompile(`\r?\n`)
		dataResponse := re.ReplaceAllString(string(dataResponseByte), "")
		dataResponse = strings.ReplaceAll(dataResponse, " ", "")
		dataResponse = strings.ReplaceAll(dataResponse, "ISO8583", "i")

		logMessage := fmt.Sprintf("\n respondStatus: %d, respondBody: %s\n", http.StatusOK, dataResponse)

		filename := fmt.Sprintf("../log/history_log/%s_%s.log", name, dateString)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		file.WriteString(logMessage)
	}
}

func IssuerLog(logMessage, name string) {
	if gin.Mode() == "debug" {
		currentTime := time.Now()
		gmtFormat := "20060102"
		dateString := currentTime.Format(gmtFormat)
		filename := fmt.Sprintf("../log/issuer_log/%s_%s.log", name, dateString)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		file.WriteString(logMessage)
	}
}
