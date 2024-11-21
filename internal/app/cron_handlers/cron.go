package cronhandlers

import (
	binrange "github.com/alfianX/hommypay_trx/databases/param/bin_range"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	transactiondata "github.com/alfianX/hommypay_trx/databases/trx/transaction_data"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	transactionDataService 	transactiondata.Service
	transactionService		transactions.Service
	hsmConfigService		hsmconfig.Service
	keyConfigService		keyconfig.Service
	binRangeService			binrange.Service
}

func NewCronjob(dbTrx *gorm.DB, dbParam *gorm.DB) CronService {
	return CronService{
		transactionDataService: transactiondata.NewService(transactiondata.NewRepo(dbTrx)),
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
		hsmConfigService: hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService: keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		binRangeService: binrange.NewService(binrange.NewRepo(dbParam)),
	}
}

func (cs *CronService) CronJob() {
	c := cron.New()

	c.AddFunc("@every 5s", func() {
		go cs.Transactions()
	})

	c.Start()

	select{}
}