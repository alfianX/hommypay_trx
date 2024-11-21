package handlers

import (
	routeconfig "github.com/alfianX/hommypay_trx/databases/param/route_config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger       	*logrus.Logger
	router       	*gin.Engine
	routeConfig 	routeconfig.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, dbParam *gorm.DB) service {
	return service{
		logger:       lg,
		router:       rtr,
		routeConfig: routeconfig.NewService(routeconfig.NewRepo(dbParam)),
	}
}