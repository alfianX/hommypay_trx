package batchupload

import (
	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/internal/batch-upload/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, lg *logrus.Logger, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB) {
	handler := handlers.NewHandler(lg, r, cnf, dbTrx, dbParam)

	r.Use(handler.MiddlewareLogger())
	r.GET("/healthz", handler.Health)
	r.POST("/batch-upload", handler.BatchUpload)
}