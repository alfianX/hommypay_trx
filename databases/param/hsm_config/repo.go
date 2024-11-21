package hsmconfig

import (
	"context"

	"gorm.io/gorm"
)

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}

func (r Repo) GetHSMIpPort(ctx context.Context) (string, string, error) {
	var hsmConfig HsmConfig
	result := r.Db.WithContext(ctx).Select("hsm_ip", "hsm_port").Find(&hsmConfig)

	return hsmConfig.HsmIp, hsmConfig.HsmPort, result.Error
}