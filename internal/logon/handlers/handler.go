package handlers

import (
	"github.com/alfianX/hommypay_trx/databases/merchant/terminals"
	hsmconfig "github.com/alfianX/hommypay_trx/databases/param/hsm_config"
	keyconfig "github.com/alfianX/hommypay_trx/databases/param/key_config"
	terminalkeys "github.com/alfianX/hommypay_trx/databases/param/terminal_keys"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	logger              *logrus.Logger
	router              *gin.Engine
	hsmConfigService    hsmconfig.Service
	keyConfigService    keyconfig.Service
	terminalKeysService terminalkeys.Service
	terminalService		terminals.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, dbParam, dbMerchant *gorm.DB) service {
	return service{
		logger:              lg,
		router:              rtr,
		hsmConfigService:    hsmconfig.NewService(hsmconfig.NewRepo(dbParam)),
		keyConfigService:    keyconfig.NewService(keyconfig.NewRepo(dbParam)),
		terminalKeysService: terminalkeys.NewService(terminalkeys.NewRepo(dbParam)),
		terminalService: terminals.NewService(terminals.NewRepo(dbMerchant)),
	}
}