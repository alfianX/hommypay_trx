package hsmconfig

import "time"

type HsmConfig struct {
	ID        int64     `json:"id"`
	HsmIp     string    `json:"hsm_ip"`
	HsmPort   string    `json:"hsm_port"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (HsmConfig) TableName() string {
	return "hsm_config"
}