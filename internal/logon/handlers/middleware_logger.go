package handlers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alfianX/hommypay_trx/helper"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
	wroteBody   bool
}

func wrapResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteBody {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.wroteBody {
		return 0, nil
	}
	i, err := rw.ResponseWriter.Write(body)
	if err != nil {
		return 0, err
	}
	rw.body = body
	return i, err
}

func (rw *responseWriter) Body() []byte {
	return rw.body
}

func (s service) MiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context){
		
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}()

		requestBody, err := helper.ReadRequestBody(c)
		if err != nil {
			helper.Respond(c, err, 0)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		helper.RestoreRequestBody(c, requestBody)

		re := regexp.MustCompile(`\r?\n`)
		reqMessage := re.ReplaceAllString(string(requestBody), "")
		reqMessage = strings.ReplaceAll(reqMessage, " ", "")

		logPrint := fmt.Sprintf("path:%s, method: %s,\n requestBody: %v", c.Request.URL.EscapedPath(), c.Request.Method, reqMessage)

		start := time.Now()
		wrapped := wrapResponseWriter(c.Writer)
		c.Writer = wrapped

		c.Next()

		logPrint = fmt.Sprintf("%s,\n respondStatus: %d, respondBody: %s", logPrint, wrapped.Status(), string(wrapped.Body()))
		
		s.logger.Infof("%s, duration: %v", logPrint, time.Since(start))

		if wrapped.Status() != 200 {
			logMessage := fmt.Sprintf("\n respondStatus: %d, respondBody: %s", wrapped.Status(), string(wrapped.Body()))

			currentTime := time.Now()
			gmtFormat := "20060102"
			dateString := currentTime.Format(gmtFormat)
			filename := fmt.Sprintf("../log/history_log/logon_%s.log", dateString)
			file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	
			if err != nil {
				helper.Respond(c, err, 0)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer file.Close()
			
			logRequest := fmt.Sprintf("%s\n", logMessage)
			file.WriteString(logRequest)
		}
	}

}