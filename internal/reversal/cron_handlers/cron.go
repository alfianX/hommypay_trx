package cronhandlers

import (
	"github.com/alfianX/hommypay_trx/configs"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	"github.com/alfianX/hommypay_trx/databases/param/issuer"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	"github.com/alfianX/hommypay_trx/databases/trx/transactions"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	config configs.Config
	reversalService reversals.Service
	hsmConfigService hsmconfig.Service
	keyConfigService keyconfig.Service
	issuerService	 issuer.Service
	transactionService transactions.Service
}

func NewCronJob(cnf configs.Config, dbTrx *gorm.DB, dbParam *gorm.DB) CronService {
	return CronService{
		config: cnf,
		reversalService: reversals.NewService(reversals.NewRepo(dbTrx)),
		hsmConfigService: hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService: keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		issuerService: issuer.NewService(issuer.NewRepo(dbParam)),
		transactionService: transactions.NewService(transactions.NewRepo(dbTrx)),
	}
}

func (cs *CronService) CronJob() {
	c := cron.New()

	c.AddFunc("@every 5s", func() {
		go cs.AutoReversal()
	})
	
	c.AddFunc("@every 30s", func() {
		go cs.SafReversal()
	})

	c.Start()

	select{}
} 