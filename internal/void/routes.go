package void

import (
	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/internal/void/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, lg *logrus.Logger, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB, dbMerchant *gorm.DB) {
	handler := handlers.NewHandler(lg, r, cnf, dbTrx, dbParam, dbMerchant)

	r.Use(handler.MiddlewareLogger())
	r.GET("/healthz", handler.Health)
	r.POST("/void", handler.Void)
}