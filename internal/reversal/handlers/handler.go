package handlers

import (
	"github.com/alfianX/hommypay_trx/configs"
	binrange "github.com/alfianX/hommypay_trx/databases/param/bin_range"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	"github.com/alfianX/hommypay_trx/databases/param/issuer"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	responsecodereversal "github.com/alfianX/hommypay_trx/databases/param/response_code_reversal"
	responsecodetrx "github.com/alfianX/hommypay_trx/databases/param/response_code_trx"
	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger             *logrus.Logger
	router             *gin.Engine
	config             configs.Config
	transactionService transactions.Service
	binRangeService    binrange.Service
	hsmConfigService   hsmconfig.Service
	keyConfigService   keyconfig.Service
	reversalService    reversals.Service
	issuerService      issuer.Service
	rcTrxService       responsecodetrx.Service
	rcReversalService  responsecodereversal.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB) service {
	return service{
		logger:             lg,
		router:             rtr,
		config:             cnf,
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
		binRangeService:    binrange.NewService(binrange.NewRepo(dbParam)),
		hsmConfigService:   hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService:   keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		reversalService:    reversals.NewService(reversals.NewRepo(dbTrx)),
		issuerService:      issuer.NewService(issuer.NewRepo(dbParam)),
		rcTrxService:       responsecodetrx.NewService(responsecodetrx.NewRepo(dbParam)),
		rcReversalService:  responsecodereversal.NewService(responsecodereversal.NewRepo(dbParam)),
	}
}
