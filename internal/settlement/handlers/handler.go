package handlers

import (
	"github.com/alfianX/hommypay_trx/configs"
	"github.com/alfianX/hommypay_trx/databases/merchant/terminals"
	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	"github.com/alfianX/hommypay_trx/databases/trx/settlement"
	settlementdetails "github.com/alfianX/hommypay_trx/databases/trx/settlement_details"
	transactionlogs "github.com/alfianX/hommypay_trx/databases/trx/transaction_logs"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger              	*logrus.Logger
	router              	*gin.Engine
	config					configs.Config
	transactionService		transactions.Service
	transactionLogService	transactionlogs.Service
	settlementService		settlement.Service
	settementDetailService	settlementdetails.Service
	terminalService			terminals.Service
	reversalService			reversals.Service
	dbTrx					*gorm.DB
	dbMerchant				*gorm.DB
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB, dbMerchant *gorm.DB) service {
	return service{
		logger:              lg,
		router:              rtr,
		config: 			 cnf,
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
		transactionLogService: transactionlogs.NewService(transactionlogs.NewRepo(dbTrx)),
		settlementService: settlement.NewService(settlement.NewRepo(dbTrx)),
		settementDetailService: settlementdetails.NewService(settlementdetails.NewRepo(dbTrx)),
		terminalService: terminals.NewService(terminals.NewRepo(dbMerchant)),
		reversalService: reversals.NewService(reversals.NewRepo(dbTrx)),
		dbTrx: dbTrx,
		dbMerchant: dbMerchant,
	}
}