package handlers

import (
	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/databases/merchant/terminals"
	"github.com/alfianX/hommypay_trx/databases/param/aid"
	binrange "github.com/alfianX/hommypay_trx/databases/param/bin_range"
	fdsconfig "github.com/alfianX/hommypay_trx/databases/param/fds_config"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	suspectlist "github.com/alfianX/hommypay_trx/databases/param/suspect_list"
	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	transactiondata "github.com/alfianX/hommypay_trx/databases/trx/transaction_data"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger              	*logrus.Logger
	router              	*gin.Engine
	config					configs.Config
	transactionDataService 	transactiondata.Service
	transactionService		transactions.Service
	binRangeService     	binrange.Service
	hsmConfigService    	hsmconfig.Service
	keyConfigService    	keyconfig.Service
	reversalService			reversals.Service
	terminalService			terminals.Service
	fdsConfigService		fdsconfig.Service
	suspectListService		suspectlist.Service
	aidService				aid.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB, dbMerchant *gorm.DB) service {
	return service{
		logger:              lg,
		router:              rtr,
		config: 			 cnf,
		transactionDataService: transactiondata.NewService(transactiondata.NewRepo(dbTrx)),
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
		binRangeService:     binrange.NewService(binrange.NewRepo(dbParam)),
		hsmConfigService:    hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService:    keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		reversalService: 	 reversals.NewService(reversals.NewRepo(dbTrx)),
		terminalService: 	 terminals.NewService(terminals.NewRepo(dbMerchant)),
		fdsConfigService: fdsconfig.NewService(fdsconfig.NewRepo(dbParam)),
		suspectListService: suspectlist.NewService(suspectlist.NewRepo(dbParam)),
		aidService: aid.NewService(aid.NewRepo(dbParam)),
	}
}