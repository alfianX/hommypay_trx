package handlers

import (
	"github.com/alfianX/hommypay_trx/configs"
	binrange "github.com/alfianX/hommypay_trx/databases/param/bin_range"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger              	*logrus.Logger
	router              	*gin.Engine
	config					configs.Config
	binRangeService     	binrange.Service
	hsmConfigService    	hsmconfig.Service
	keyConfigService    	keyconfig.Service
	transactionService		transactions.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB) service {
	return service{
		logger:              lg,
		router:              rtr,
		config: 			 cnf,
		binRangeService:     binrange.NewService(binrange.NewRepo(dbParam)),
		hsmConfigService:    hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService:    keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
	}
}