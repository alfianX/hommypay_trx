package app

import (
	"github.com/alfianX/hommypay_trx/internal/app/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, lg *logrus.Logger, dbParam *gorm.DB) {
	handler := handlers.NewHandler(lg, r, dbParam)

	r.Use(handler.MiddlewareLogger())
	s := r.Group("/api/uplink/v1")
	s.GET("/healthz", handler.Health)
	s.POST("/*action", handler.Actions)
	// s.GET("/createKeyRSA", handler.CreateKeyRSA)
}