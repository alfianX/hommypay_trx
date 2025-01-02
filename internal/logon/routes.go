package logon

import (
	"github.com/alfianX/hommypay_trx/internal/logon/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, lg *logrus.Logger, dbParam, dbMerchant *gorm.DB) {
	handler := handlers.NewHandler(lg, r, dbParam, dbMerchant)

	r.Use(handler.MiddlewareLogger())
	r.GET("/healthz", handler.Health)
	r.POST("/logon", handler.Logon)
}